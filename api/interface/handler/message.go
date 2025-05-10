package handler

import (
	"chat_back/usecase"
	"encoding/json"
	"fmt"
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

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
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

type message struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
}

type messageURI struct {
	ChannelID string `uri:"channelID" binding:"required,uuid"`
}

type MessageHandler interface {
	HandleMessageByID(ctx *gin.Context)
	HandleMessageInChannel(ctx *gin.Context)
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
	// send      chan []byte
	send chan message
}

type Hub struct {
	clients map[string]map[*Client]struct{}
	// clients    map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
	// broadcast  chan []byte
	broadcast chan message
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			fmt.Println("register")
			if _, ok := h.clients[client.channelID]; !ok {
				h.clients[client.channelID] = make(map[*Client]struct{})
			}
			h.clients[client.channelID][client] = struct{}{}

		case client := <-h.unregister:
			fmt.Println("unregister")
			if _, ok := h.clients[client.channelID][client]; ok {
				delete(h.clients[client.channelID], client)
				close(client.send)
				fmt.Printf("Websocket in %s closed", client.channelID)
				if len(h.clients[client.channelID]) <= 0 {
					delete(h.clients, client.channelID)
					fmt.Printf("All websocket in %s closed", client.channelID)
				}
			}

		case msg := <-h.broadcast:
			fmt.Println("broadcast")
			fmt.Println("msg", msg)
			for client := range h.clients[msg.ChannelID] {
				select {
				case client.send <- msg:
					fmt.Println("msg sent to ", client.channelID, client.userID)
				default:
					close(client.send)
					delete(h.clients[client.channelID], client)
				}
			}
		}

	}
}

func (c *Client) readPump(uc usecase.MessageUsecase, user *clerk.User) {
	fmt.Println("read pump")
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
		fmt.Println("reading")
		_, rawMsg, err := c.conn.ReadMessage()
		fmt.Println(string(rawMsg), err)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("unexpected close error:", err)
			}
			break
		}

		var msg message
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			fmt.Println("fail to unmarshal:", rawMsg)
			break
		}
		fmt.Println("read", msg)

		insertedMsg, err := uc.Insert(msg.ChannelID, user.ID, msg.Content)
		if err != nil {
			fmt.Println("message insertion error", err)
			break
		}

		msg.ID = insertedMsg.ID
		msg.UserName = *user.Username

		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	fmt.Println("write pump")

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		fmt.Println("writing", c.channelID)
		select {
		case msg, ok := <-c.send:
			fmt.Println("get msg <- send")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			parsed_msg, err := json.Marshal(msg)
			if err != nil {
				fmt.Println("fail to marshal:", msg)
			}

			c.conn.WriteMessage(websocket.TextMessage, parsed_msg)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func NewMessageHandler(messageUseCase usecase.MessageUsecase, authorizationUseCase usecase.AuthorizationUsecase) MessageHandler {
	fmt.Println("NewMessageHandler")

	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	hub := Hub{
		clients: map[string]map[*Client]struct{}{},
		// clients:    make(map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan message),
		// broadcast: make(chan []byte),
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
	fmt.Println("HandleMessageWebSocket")

	var messageURI messageURI
	if err := ctx.ShouldBindUri(&messageURI); err != nil {
		fmt.Println("Error binding URI:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	conn, err := mh.wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("websocket upgrade err:", err)
		return
	}

	client := &Client{hub: mh.hub, conn: conn, channelID: messageURI.ChannelID, send: make(chan message, 256)}
	client.hub.register <- client

	msgType, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("authorization msg", err)
		conn.Close()
		return
	}
	if msgType != websocket.TextMessage {
		fmt.Println("authorization msg", msgType)
		conn.Close()
		return
	}

	var authorizationMessage authorizationMessage
	if err := json.Unmarshal(msg, &authorizationMessage); err != nil {
		fmt.Println("Error unmarshalling message:", err)
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
		fmt.Println("Error checking permission:", err)
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

	var messageURI messageURI
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

func (mh messageHandler) HandleMessageInChannel(ctx *gin.Context) {
	fmt.Println("HandleMessageInChannel")

	var messageURI messageURI
	if err := ctx.ShouldBindUri(&messageURI); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	fmt.Println("channelID", messageURI.ChannelID)

	messages, err := mh.messageUseCase.GetAllInChannel(messageURI.ChannelID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get messages"})
		return
	}

	var parsed_messages []message
	for _, msg := range messages {
		parsed_messages = append(parsed_messages, message{
			ID:        msg.ID,
			UserID:    msg.UserID,
			UserName:  msg.UserName,
			ChannelID: msg.ChannelID,
			Content:   msg.Content,
		})
	}
	fmt.Println("parsed_messages", parsed_messages)

	ctx.JSON(http.StatusOK, parsed_messages)
}
