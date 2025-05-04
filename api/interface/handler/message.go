package handler

import (
	"chat_back/usecase"
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type message struct {
	userID    string
	channelID string
	content   []byte
}

type messageURI struct {
	ChannelID string `uri:"channelID" binding:"required,uuid"`
}

type authorizationMessage struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	Token     string `json:"token"`
}

type MessageHandler interface {
	HandleMessageByID(ctx *gin.Context)
	HandleMessageInChannel(ctx *gin.Context)
	HandleMessageWebSocket(ctx *gin.Context)
}

type messageHandler struct {
	db                   *gorm.DB
	wsUpgrader           *websocket.Upgrader
	hub                  *Hub
	messageUseCase       usecase.MessageUsecase
	authorizationUseCase usecase.AuthorizationUsecase
}

type Client struct {
	conn      *websocket.Conn
	userID    string
	channelID string
}

type Hub struct {
	conns      *map[string]map[*websocket.Conn]struct{}
	register   chan *Client
	unregister chan *Client
	broadcast  chan message
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			if _, ok := (*h.conns)[client.channelID]; !ok {
				(*h.conns)[client.channelID] = make(map[*websocket.Conn]struct{})
			}
			(*h.conns)[client.channelID][client.conn] = struct{}{}
		case client := <-h.unregister:
			if _, ok := (*h.conns)[client.channelID]; ok {
				if _, ok := (*h.conns)[client.channelID][client.conn]; ok {
					delete((*h.conns)[client.channelID], client.conn)
					if len((*h.conns)[client.channelID]) == 0 {
						delete(*h.conns, client.channelID)
					}
				}
			}
		case msg := <-h.broadcast:
			for conn, _ := range (*h.conns)[msg.channelID] {
				err := conn.WriteMessage(websocket.TextMessage, msg.content)
				if err != nil {
					fmt.Println("Error writing message:", err)
					conn.Close()
					delete((*h.conns)[msg.channelID], conn)
				}
			}

		}
	}
}

func NewMessageHandler(db *gorm.DB, messageUseCase usecase.MessageUsecase, authorizationUseCase usecase.AuthorizationUsecase) MessageHandler {
	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	hub := Hub{
		conns:      &map[string]map[*websocket.Conn]struct{}{},
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan message),
	}
	go hub.run()

	return &messageHandler{
		db:                   db,
		wsUpgrader:           &wsUpgrader,
		hub:                  &hub,
		messageUseCase:       messageUseCase,
		authorizationUseCase: authorizationUseCase,
	}
}

func waitForMessage(uc usecase.MessageUsecase, db *gorm.DB, messageURI messageURI, conn *websocket.Conn, broadcast chan message) {
	for {
		msgType, msg, err := conn.ReadMessage()
		fmt.Println("read a message", msgType, string(msg))
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		if msgType == websocket.TextMessage {
			fmt.Println("broadcast")
			uc.Insert(db, messageURI.ChannelID, "userID", string(msg))
			broadcast <- message{
				userID:    "userID",
				channelID: messageURI.ChannelID,
				content:   msg,
			}
		}
	}
}

func (mh messageHandler) HandleMessageWebSocket(ctx *gin.Context) {
	fmt.Println("HandleMessageWebSocket")
	var messageURI messageURI
	if err := ctx.ShouldBindUri(&messageURI); err != nil {
		fmt.Println("Error binding URI:", err)
		ctx.String(http.StatusBadRequest, "Bad request")
		return
	}

	conn, err := mh.wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		ctx.String(http.StatusInternalServerError, "Failed to upgrade connection")
		return
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println(err)
		ctx.String(http.StatusInternalServerError, "Failed to read message")
		return
	}
	fmt.Println("first", string(msg))
	var authorizationMessage authorizationMessage
	if err := json.Unmarshal(msg, &authorizationMessage); err != nil {
		fmt.Println("Error unmarshalling message:", err)
		ctx.String(http.StatusBadRequest, "Invalid message format")
		return
	}
	fmt.Println("authorizationMessage", authorizationMessage)
	if authorizationMessage.ChannelID != messageURI.ChannelID {
		fmt.Println("Channel ID mismatch")
		ctx.String(http.StatusBadRequest, "Channel ID mismatch")
		return
	}

	fmt.Println("hi")
	fmt.Println(authorizationMessage.UserID)
	fmt.Println(authorizationMessage.Token)
	fmt.Println(authorizationMessage.ChannelID)
	fmt.Println("hi")

	user, err := mh.authorizationUseCase.CheckPermission(mh.db, authorizationMessage.UserID, authorizationMessage.ChannelID, authorizationMessage.Token)
	if err != nil {
		fmt.Println("Error checking permission:", err)
		ctx.String(http.StatusInternalServerError, "Failed to check permission")
		return
	}
	if user == nil {
		fmt.Println("User not found or no permission")
		ctx.String(http.StatusForbidden, "No permission")
		return
	}

	mh.hub.register <- &Client{
		conn:      conn,
		userID:    authorizationMessage.UserID,
		channelID: messageURI.ChannelID,
	}

	go waitForMessage(mh.messageUseCase, mh.db, messageURI, conn, mh.hub.broadcast)
}

func (mh messageHandler) HandleMessageByID(ctx *gin.Context) {
	var messageURI messageURI
	if err := ctx.ShouldBindUri(&messageURI); err != nil {
		ctx.String(http.StatusBadRequest, "Invalid request")
		return
	}

	message, err := mh.messageUseCase.GetByID(mh.db, messageURI.ChannelID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to get message")
		return
	}
	ctx.JSON(http.StatusOK, message)
}

func (mh messageHandler) HandleMessageInChannel(ctx *gin.Context) {
	var messageURI messageURI
	fmt.Println(ctx)
	if err := ctx.ShouldBindUri(&messageURI); err != nil {
		ctx.String(http.StatusBadRequest, "Invalid request")
		return
	}
	fmt.Println("channelID", messageURI.ChannelID)

	messages, err := mh.messageUseCase.GetAllInChannel(mh.db, messageURI.ChannelID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to get messages")
		return
	}
	ctx.JSON(http.StatusOK, messages)
}
