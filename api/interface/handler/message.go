package handler

import (
	"bytes"
	"chat_back/usecase"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
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
	db                   *gorm.DB
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
	send      chan []byte
}

type Hub struct {
	// clients    *map[string]map[*Client]struct{}
	clients    map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

// func (h *Hub) run() {
// 	for {
// 		select {
// 		case client := <-h.register:
// 			if _, ok := (*h.clients)[client.channelID]; !ok {
// 				(*h.clients)[client.channelID] = make(map[*Client]struct{})
// 			}
// 			(*h.clients)[client.channelID][client] = struct{}{}
// 			fmt.Println("a client registered:", client.channelID)

// 		case client := <-h.unregister:
// 			if _, ok := (*h.clients)[client.channelID]; ok {
// 				if _, ok := (*h.clients)[client.channelID][client]; ok {
// 					fmt.Println("closed: ", client.channelID)
// 					client.conn.Close()
// 					delete((*h.clients)[client.channelID], client)
// 					if len((*h.clients)[client.channelID]) == 0 {
// 						fmt.Println("all closed:", client.channelID)
// 						delete(*h.clients, client.channelID)
// 					}
// 				}
// 			}

// 		case msg := <-h.broadcast:
// 			parsedMessage, err := json.Marshal(msg)
// 			if err != nil {
// 				fmt.Println("Error marshalling message:", err)
// 				continue
// 			}
// 			for client, _ := range (*h.clients)[msg.ChannelID] {
// 				fmt.Println("broadcast:", msg.ChannelID)

// 				if client.channelID != msg.ChannelID {
// 					fmt.Println("channel id mismatch:", client.channelID)
// 				}

// 				err = client.conn.WriteMessage(websocket.TextMessage, parsedMessage)
// 				if err != nil {
// 					fmt.Println("Error writing message:", err)
// 					h.unregister <- client
// 				}
// 			}

// 		}
// 	}
// }

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			fmt.Println("register")
			h.clients[client] = struct{}{}
		case client := <-h.unregister:
			fmt.Println("unregister")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case msg := <-h.broadcast:
			fmt.Println("broadcast")
			fmt.Println(msg)
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}

	}
}

func (c *Client) readPump() {
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
		_, msg, err := c.conn.ReadMessage()
		fmt.Println(msg, err)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("unexpected close error:", err)
			}
			break
		}

		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))

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
		fmt.Println("writing")
		select {
		case msg, ok := <-c.send:
			fmt.Println("get msg <- send")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			writer.Write(msg)

			for i := 0; i < len(c.send); i++ {
				writer.Write(newline)
				writer.Write(<-c.send)
			}

			if err := writer.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func NewMessageHandler(db *gorm.DB, messageUseCase usecase.MessageUsecase, authorizationUseCase usecase.AuthorizationUsecase) MessageHandler {
	fmt.Println("NewMessageHandler")

	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	hub := Hub{
		// clients:    &map[string]map[*Client]struct{}{},
		clients:    make(map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		// broadcast:  make(chan message),
		broadcast: make(chan []byte),
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

// func waitForMessage(uc usecase.MessageUsecase, db *gorm.DB, user *clerk.User, messageURI messageURI, conn *websocket.Conn, broadcast chan message) {
// 	fmt.Println("waitForMessage")

// 	for {
// 		msgType, msg, err := conn.ReadMessage()
// 		fmt.Println("read a message", msgType, string(msg))
// 		if err != nil {
// 			fmt.Println("Error reading message:", err)
// 			break
// 		}
// 		if msgType != websocket.TextMessage {
// 			fmt.Println("Message type mismatch:", msgType)
// 		}

// 		var sentMessage sentMessage
// 		if err := json.Unmarshal(msg, &sentMessage); err != nil {
// 			fmt.Println("Error unmarshalling message:", err)
// 			break
// 		}
// 		fmt.Println("receiveMessage", sentMessage)

// 		if sentMessage.Action != Send {
// 			fmt.Println("Action mismatch:", sentMessage.Action)
// 			break
// 		}
// 		if sentMessage.ChannelID != messageURI.ChannelID {
// 			fmt.Println("Channel ID mismatch")
// 			break
// 		}

// 		fmt.Println("broadcast")
// 		insertedMessage, err := uc.Insert(db, sentMessage.ChannelID, user.ID, sentMessage.Content)
// 		if err != nil {
// 			fmt.Println(err)
// 			continue
// 		}

// 		broadcast <- message{
// 			ID:        insertedMessage.ID,
// 			UserID:    insertedMessage.UserID,
// 			UserName:  insertedMessage.UserName,
// 			ChannelID: insertedMessage.ChannelID,
// 			Content:   insertedMessage.Content,
// 		}
// 	}
// }

func (mh messageHandler) HandleMessageWebSocket(ctx *gin.Context) {
	// fmt.Println("HandleMessageWebSocket")

	// var messageURI messageURI
	// if err := ctx.ShouldBindUri(&messageURI); err != nil {
	// 	fmt.Println("Error binding URI:", err)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	// 	return
	// }

	// conn, err := mh.wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade connection"})
	// 	return
	// }

	conn, err := mh.wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("websocket upgrade err:", err)
		return
	}
	client := &Client{hub: mh.hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()

	// _, msg, err := conn.ReadMessage()
	// if err != nil {
	// 	fmt.Println(err)
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read message"})
	// 	return
	// }

	// var authorizationMessage authorizationMessage
	// if err := json.Unmarshal(msg, &authorizationMessage); err != nil {
	// 	fmt.Println("Error unmarshalling message:", err)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid message format"})
	// 	return
	// }
	// fmt.Println("authorizationMessage", authorizationMessage)
	// if authorizationMessage.ChannelID != messageURI.ChannelID {
	// 	fmt.Println("Channel ID mismatch")
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "channel ID mismatch"})
	// 	return
	// }

	// fmt.Println("hi")
	// fmt.Println(authorizationMessage.ChannelID)
	// fmt.Println("hi")

	// user, err := mh.authorizationUseCase.CheckPermission(mh.db, authorizationMessage.ChannelID, authorizationMessage.Token)
	// if err != nil {
	// 	fmt.Println("Error checking permission:", err)
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permission"})
	// 	return
	// }
	// if user == nil {
	// 	fmt.Println("User not found or no permission")
	// 	ctx.JSON(http.StatusForbidden, gin.H{"error": "no permission"})
	// 	return
	// }

	// mh.hub.register <- &Client{
	// 	hub:       mh.hub,
	// 	conn:      conn,
	// 	userID:    user.ID,
	// 	channelID: messageURI.ChannelID,
	// 	send:      make(chan []byte, 256),
	// }

	// go waitForMessage(mh.messageUseCase, mh.db, user, messageURI, conn, mh.hub.broadcast)
}

func (mh messageHandler) HandleMessageByID(ctx *gin.Context) {
	fmt.Println("HandleMessageByID")

	var messageURI messageURI
	if err := ctx.ShouldBindUri(&messageURI); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	message, err := mh.messageUseCase.GetByID(mh.db, messageURI.ChannelID)
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

	messages, err := mh.messageUseCase.GetAllInChannel(mh.db, messageURI.ChannelID)
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
