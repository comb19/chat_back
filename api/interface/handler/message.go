package handler

import (
	"chat_back/interface/types"
	"chat_back/usecase"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	Authorization = "authorization"
	Send          = "send"
)

const (
	writeWait      = 10 * time.Second
	pingPeriod     = 1 * time.Second
	pongWait       = 60 * time.Second
	maxMessageSize = 4096
)

type sentMessage struct {
	Action    string `json:"action"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
}

type authorizationMessage struct {
	Action    string `json:"action"`
	ChannelID string `json:"channel_id"`
	Token     string `json:"token"`
}

type MessageHandler interface {
	HandleMessageByID(ctx *gin.Context)
	HandleMessageWebSocket(ctx *gin.Context)
}

type messageHandler struct {
	wsUpgrader           *websocket.Upgrader
	hub                  *Hub
	messageUseCase       usecase.MessageUsecase
	authorizationUseCase usecase.AuthorizationUsecase
}

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	userID    string
	channelID string
	send      chan types.Message
}

type Hub struct {
	clients    map[string]map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
	broadcast  chan types.Message
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			slog.Debug("register a websocket connection")

			if _, ok := h.clients[client.channelID]; !ok {
				h.clients[client.channelID] = make(map[*Client]struct{})
			}
			h.clients[client.channelID][client] = struct{}{}

		case client := <-h.unregister:
			slog.Debug("unregister a websocket connection")

			if _, ok := h.clients[client.channelID][client]; ok {
				delete(h.clients[client.channelID], client)
				close(client.send)
				slog.Debug(fmt.Sprintf("Websocket in %s closed", client.channelID))

				if len(h.clients[client.channelID]) <= 0 {
					delete(h.clients, client.channelID)
					slog.Debug(fmt.Sprintf("All websocket in %s closed", client.channelID))
				}
			}

		case msg := <-h.broadcast:
			slog.Debug("broadcast")
			slog.Debug(fmt.Sprintf("msg: %s", msg))

			for client := range h.clients[msg.ChannelID] {
				select {
				case client.send <- msg:
					slog.Debug(fmt.Sprintf("msg sent to %s by %s", client.channelID, client.userID))
				default:
					close(client.send)
					delete(h.clients[client.channelID], client)
				}
			}
		}

	}
}

func (c *Client) readPump(uc usecase.MessageUsecase, user *clerk.User) {
	slog.Debug("readPump")

	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Debug(err.Error())
			}
			break
		}

		var msg types.Message
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			slog.Debug(err.Error())
			break
		}
		fmt.Println("read", msg)

		insertedMsg, err := uc.Insert(msg.ChannelID, user.ID, msg.Content)
		if err != nil {
			slog.Debug(err.Error())
			break
		}

		msg.ID = insertedMsg.ID
		msg.UserName = *user.Username
		msg.CreatedAt = insertedMsg.CreatedAt
		msg.UserID = insertedMsg.UserID

		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	slog.Debug("writePump")

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			parsed_msg, err := json.Marshal(msg)
			if err != nil {
				slog.Debug(err.Error())
				continue
			}

			c.conn.WriteMessage(websocket.TextMessage, parsed_msg)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Debug(err.Error())
				return
			}
		}
	}
}

func NewMessageHandler(messageUseCase usecase.MessageUsecase, authorizationUseCase usecase.AuthorizationUsecase) MessageHandler {
	slog.Debug("NewMessageHandler")

	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	hub := Hub{
		clients:    map[string]map[*Client]struct{}{},
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan types.Message),
	}
	go hub.run()

	return &messageHandler{
		wsUpgrader:           &wsUpgrader,
		hub:                  &hub,
		messageUseCase:       messageUseCase,
		authorizationUseCase: authorizationUseCase,
	}
}

func (mh messageHandler) HandleMessageWebSocket(ctx *gin.Context) {
	slog.DebugContext(ctx, "HandleMessageWebSocket")

	var messageURI types.MessageURI
	if err := ctx.ShouldBindUri(&messageURI); err != nil {
		slog.ErrorContext(ctx, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	conn, err := mh.wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		fmt.Println("websocket upgrade err:", err)
		return
	}

	client := &Client{hub: mh.hub, conn: conn, channelID: messageURI.ChannelID, send: make(chan types.Message, 256)}
	client.hub.register <- client

	msgType, msg, err := conn.ReadMessage()
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		conn.Close()
		return
	}
	if msgType != websocket.TextMessage {
		slog.ErrorContext(ctx, fmt.Sprintf("authorization msg: %d", msgType))
		conn.Close()
		return
	}

	var authorizationMessage authorizationMessage
	if err := json.Unmarshal(msg, &authorizationMessage); err != nil {
		slog.ErrorContext(ctx, err.Error())
		conn.Close()
		return
	}
	fmt.Println("authorizationMessage", authorizationMessage)
	if authorizationMessage.ChannelID != messageURI.ChannelID {
		fmt.Println("Channel ID mismatch")
		conn.Close()
		return
	}

	user, err := mh.authorizationUseCase.CheckPermission(authorizationMessage.ChannelID, authorizationMessage.Token)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		conn.Close()
		return
	}
	if user == nil {
		fmt.Println("User not found or no permission")
		conn.Close()
		return
	}

	go client.writePump()
	go client.readPump(mh.messageUseCase, user)

}

func (mh messageHandler) HandleMessageByID(ctx *gin.Context) {
	fmt.Println("HandleMessageByID")

	var messageURI types.MessageURI
	if err := ctx.ShouldBindUri(&messageURI); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	message, err := mh.messageUseCase.GetByID(messageURI.ChannelID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get message"})
		return
	}
	ctx.JSON(http.StatusOK, message)
}
