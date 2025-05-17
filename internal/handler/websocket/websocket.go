package GRWebsocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/azan-boss/posty/internal/handler/auth"

	"github.com/azan-boss/posty/internal/storage"
	"github.com/azan-boss/posty/internal/types"
	"github.com/azan-boss/posty/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
)

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"sync"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// )

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// )

// // upgrader is used to upgrade HTTP connections to WebSocket connections
// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true // TODO: implement origin check
// 	},
// }

// var clients = make(map[*websocket.Conn]string)

// func HandleWebSocket(c *gin.Context)  {
// 	// upgrade the connection to a WebSocket connection
// 	conn ,err :=upgrader.Upgrade(c.Writer,c.Request,nil)

// 	if err!=nil{
// 		log.Fatal(err)
// 	}

// 	defer conn.Close()
// 	_,message,err := conn.ReadMessage()
// 	if err!=nil{
// 		log.Fatal(err)
// 	}
// 	json.Unmarshal()
// 	log.Println(string(message))
// 	clients[conn] =string(message)
// 	broadcast(string(message) + " is online")

// 	for {
// 		_, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			break
// 		}
// 		message := string(message) + ": " + string(msg)
// 		broadcast(message)
// 	}

// 	delete(clients,conn)
// 	broadcast(string(message) + " is offline")
// 	// for client, message := range clients {
// 	// 	log.Println("Client:", client, "Message:", message)
// 	// }

// }

// func broadcast(message string) {
// 	for conn := range clients {
// 		conn.WriteMessage(websocket.TextMessage, []byte(message))
// 	}
// }

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true // Adjust for production
// 	},
// }

// type Client struct {
// 	conn     *websocket.Conn
// 	username string
// 	send     chan Message
// }

// type Message struct {
// 	Username  string    `json:"username"`
// 	Content   string    `json:"content"`
// 	Type      string    `json:"type"` // "message", "notification", "typing"
// 	Timestamp time.Time `json:"timestamp"`
// }

// var (
// 	clients   = make(map[*Client]bool)
// 	clientsMu sync.Mutex // To protect concurrent access to clients map
// 	broadcast = make(chan Message)
// )

// func HandleWebSocket(c *gin.Context) {
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println("Upgrade error:", err)
// 		return
// 	}

// 	// First message should be the username
// 	_, msg, err := conn.ReadMessage()
// 	if err != nil {
// 		log.Println("Initial read error:", err)
// 		conn.Close()
// 		return
// 	}

// 	var initMsg struct {
// 		Username string `json:"username"`
// 	}
// 	if err := json.Unmarshal(msg, &initMsg); err != nil {
// 		log.Println("JSON parse error:", err)
// 		conn.Close()
// 		return
// 	}

// 	client := &Client{
// 		conn:     conn,
// 		username: initMsg.Username,
// 		send:     make(chan Message, 256),
// 	}

// 	// Register client
// 	clientsMu.Lock()
// 	clients[client] = true
// 	clientsMu.Unlock()

// 	// Notify others
// 	broadcast <- Message{
// 		Username:  client.username,
// 		Content:   "has joined the chat",
// 		Type:      "notification",
// 		Timestamp: time.Now(),
// 	}

// 	// Start goroutines for read/write
// 	go client.readPump()
// 	go client.writePump()
// }

// func (c *Client) readPump() {
// 	defer func() {
// 		// Cleanup on exit
// 		clientsMu.Lock()
// 		delete(clients, c)
// 		clientsMu.Unlock()
// 		c.conn.Close()

// 		// Notify others
// 		broadcast <- Message{
// 			Username:  c.username,
// 			Content:   "has left the chat",
// 			Type:      "notification",
// 			Timestamp: time.Now(),
// 		}
// 	}()

// 	for {
// 		_, msg, err := c.conn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}

// 		var incoming struct {
// 			Content string `json:"content"`
// 			Type    string `json:"type"`
// 		}
// 		if err := json.Unmarshal(msg, &incoming); err != nil {
// 			log.Println("JSON parse error:", err)
// 			continue
// 		}

// 		// Handle different message types
// 		switch incoming.Type {
// 		case "message":
// 			broadcast <- Message{
// 				Username:  c.username,
// 				Content:   incoming.Content,
// 				Type:      "message",
// 				Timestamp: time.Now(),
// 			}
// 		case "typing":
// 			broadcast <- Message{
// 				Username:  c.username,
// 				Content:   "is typing...",
// 				Type:      "typing",
// 				Timestamp: time.Now(),
// 			}
// 		}
// 	}
// }

// func (c *Client) writePump() {
// 	defer c.conn.Close()

// 	for {
// 		message, ok := <-c.send
// 		if !ok {
// 			// Channel closed
// 			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
// 			return
// 		}

// 			msgBytes, err := json.Marshal(message)
// 			if err != nil {
// 				log.Println("JSON marshal error:", err)
// 				return
// 			}

// 			if err := c.conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
// 				log.Println("Write error:", err)
// 				return
// 			}

// 	}
// }

// func HandleBroadcast() {
// 	for {
// 		msg := <-broadcast

// 		clientsMu.Lock()
// 		for client := range clients {
// 			select {
// 			case client.send <- msg:
// 			default:
// 				// Couldn't send, close connection
// 				close(client.send)
// 				delete(clients, client)
// 			}
// 		}
// 		clientsMu.Unlock()
// 	}
// }

// type Client struct{
// 	conn *websocket.Conn
// 	user types.User
// 	roomId string
// 	send chan types.Message
// }

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true // TODO: implement origin check
// 	},
// }

// var rooms = make(map[string]map[*Client]bool)
// var ClientMux sync.Mutex

// func HandleWebSocket(storage storage.Storage) (gin.HandlerFunc) {
// 	return func (c *gin.Context)  {
// 		ChatRoomId := c.Param("roomId")
// 		token:= c.Query("token")

// 		if token==""{
// 			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("unauthorized")))
// 			return
// 		}
// 		// fmt.Println(token)
// 		claims, err := auth.VerifyJWT(token)
// 		if err != nil {
// 			slog.Error("Invalid token", "error", err)
// 			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("unauthorized")))
// 			return
// 		}
// 		userIdString := strconv.Itoa(int(claims.UserId))
// 		chatroom, error := storage.GetUseChatRoom(ChatRoomId,userIdString)
// 		slog.Info("Chatroom found", "chatroom", chatroom, "user_id", userIdString)
// 	if error!=nil{
// 		slog.Error("Chatroom not found", "error", error)
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Chatroom not found"})
// 		return
// 	}

// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		slog.Error("Failed to upgrade connection", "error", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
// 		return
// 	}
// 	defer conn.Close()

// 	 user, _ := storage.GetUser(uint(claims.UserId))
// 	 client := &Client{
// 		 conn: conn,
// 		user: user,
// 		roomId: ChatRoomId,
// 		send: make(chan types.Message),
// 	}

// 	ClientMux.Lock()
// 	if _, ok := rooms[ChatRoomId]; !ok {
// 		rooms[ChatRoomId] = make(map[*Client]bool)
// 	}
// 	rooms[ChatRoomId][client] = true
// 	ClientMux.Unlock()

// 	go client.writePump(storage,user)
// 	go client.readPump(storage,user)

// }}

// func (c *Client) readPump(storage storage.Storage, user types.User) {
// 	defer func() {
// 		ClientMux.Lock()
// 		c.conn.Close()
// 		storage.UpdateStatus(&user)
// 		delete(rooms[c.roomId], c)
// 		ClientMux.Unlock()
// 	}()

// 	convId, _ := strconv.Atoi(c.roomId)
// 	for {
// 		_, message, err := c.conn.ReadMessage()
// 		if err != nil {
// 			slog.Error("Failed to read message", "error", err)
// 			break // Exit the loop on error
// 		}

// 		var msg types.Message
// 		if err := json.Unmarshal(message, &msg); err != nil {
// 			slog.Error("JSON parse error", "error", err)
// 			continue // Skip to the next iteration
// 		}

// 		switch msg.Type {
// 		case "typing":
// 			Brodcast(c, types.Message{
// 				Content: fmt.Sprintf("%s is typing", c.user.Username),
// 				Type:    "typing",
// 				UserID:  c.user.ID,
// 				ChatRoomId: uint(convId),
// 			}, storage, user)
// 		default:
// 			Brodcast(c, types.Message{
// 				Content: msg.Content,
// 				Type:    "message",
// 				UserID:  c.user.ID,
// 				ChatRoomId: uint(convId),
// 			}, storage, user)
// 		}
// 	}
// }

// func (c *Client)writePump(storage storage.Storage ,user types.User)  {
// 	defer func(){
// 		ClientMux.Lock()
// 		c.conn.Close()
// 		storage.UpdateStatus(&user)
// 		delete(rooms[c.roomId],c)
// 		ClientMux.Unlock()

// 		}()

// 		 for {
// 			msg, ok := <-c.send
// 			if !ok {
// 				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				continue
// 			}

// 			messageByte ,err:=json.Marshal(msg)
// 			if err!=nil{
// 				slog.Error("Failed to marshal message", "error", err)
// 				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				continue

// 			}
// 			c.conn.WriteMessage(websocket.TextMessage,messageByte)
// 		 }
// }

// func Brodcast(c *Client,message types.Message,storage storage.Storage,user types.User)  {
// 	ClientMux.Lock()
// 	defer ClientMux.Unlock()
// 	for client :=range rooms[c.roomId]{
// 		select{

// 		case client.send <-message:
// 		default:
// 			c.conn.Close()
// 			ClientMux.Lock()
// 			delete(rooms[c.roomId],c)
// 			c.conn.Close()
// 			user.Status = "offline"
// 			storage.UpdateStatus(&user)
// 			ClientMux.Unlock()

// 		}
// 	}
// }








type  Client struct{
	coon *websocket.Conn
	user types.User
	roomId string
	Send chan types.Message
}
var done = make(chan bool,1)

var rooms = make(map[string]map[*Client]bool)
var ClientMux sync.Mutex 



var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: implement origin check
	},
}


func HandleWebSocket(storage storage.Storage)  gin.HandlerFunc {
	return func(c *gin.Context) {
		
		// 1-Handler authentication  
		token:=c.Query("token")
		if token==""{
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("unauthorized")))
			slog.Info("token does not provided")
			return
		}
		
		claims, err := auth.VerifyJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("unauthorized")))
			slog.Info("jwt token verification failed")
			return
		}
		slog.Info("Login successfully %s",claims.UserId)
		

		// Getting information of user and chatroom
		user, err := storage.GetUser(uint(claims.UserId))
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("unauthorized")))
			slog.Error("User not found", "error", err)
			return
		}

		chatroom, err := storage.FindUserChatRoomOrJoin(c.Param("roomId"),strconv.Itoa(int(claims.UserId)))
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("unauthorized")))
			slog.Error("Chatroom not found", "error", err)
			return
		}
		slog.Info("Chatroom found", "chatroom", chatroom, "user_id", claims.UserId)
		
		
		// 4-Upgrade connection htto to websocket
		coon, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to upgrade connection")))
			slog.Error("Failed to upgrade connection", "error", err)
			return
		}



		// 5-Create client of current user 
		client := &Client{
			coon: coon,
			user: user,
			roomId: strconv.Itoa(int(chatroom.ID)),
			Send: make(chan types.Message),
		}
		ClientMux.Lock()
		if _, ok := rooms[client.roomId]; !ok {
			rooms[client.roomId] = make(map[*Client]bool)
		}
		rooms[client.roomId][client] = true
		ClientMux.Unlock()
		messages, err := storage.GetMessageByChatRoomId(client.roomId)
		if err != nil {
			slog.Error("Error while fetching message history", "error", err)
		} else {
			// Send message history as a special message type
			historyMsg := map[string]interface{}{
				"type":     "history",
				"messages": messages,
			}
			if err := client.coon.WriteJSON(historyMsg); err != nil {
				slog.Error("Failed to send message history", "error", err)
			}
			slog.Info("Message history sent successfully")
		}

		//  6.Starting goroutines  to handle out the  send and receive messages 
	go client.readPump(storage)
	go client.writePump(storage)

	}
}



func (client *Client)readPump(storage storage.Storage){
	defer func ()  {
		ClientMux.Lock()
		delete(rooms[client.roomId],client)
		client.coon.Close()
		ClientMux.Unlock()
		done <- true
		}()

	for {
		var incoming struct{
			Content string
			Type string
		}

		_,msg, err:=client.coon.ReadMessage()

		 if err!=nil{
			slog.Error("error while fetching the message")
			break
		 }
		
		 err=json.Unmarshal(msg,&incoming)

		 if err!=nil{
			slog.Error("Failed to pras message")
			continue
		 }

		 roomId,_:=strconv.Atoi(client.roomId)

		 
		 switch incoming.Type {
		 case "typing":
			PublishQueue(types.Message{
				Content: incoming.Content,
				Type:    "typing",
				UserID:  client.user.ID,
				ChatRoomId: uint(roomId),
			})
			
		 default:
			msg:=types.Message{
				Content: incoming.Content,
				Type:    "message",
				UserID:  client.user.ID,
				ChatRoomId: uint(roomId),
			}
			PublishQueue(msg)
			client.Send <- msg
		 	
		 }
	}
}

func (client *Client) writePump(storage storage.Storage) {
    defer func() {
        ClientMux.Lock()
        delete(rooms[client.roomId], client)
        client.coon.Close()
        ClientMux.Unlock()
    }()

    for {
        select {
        case <-done:
            slog.Info("WritePump: done channel closed, exiting writePump")
            return
        case msg, ok := <-client.Send:
            if !ok {
                slog.Warn("Send channel closed unexpectedly")
                return
            }

            msgBytes, err := json.Marshal(msg)
            if err != nil {
                slog.Error("Error while converting to JSON:", "error", err)
                continue
            }

            if err := client.coon.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
                slog.Error("Error while writing message:", "error", err)
                continue
            }
			ChatRoomId,_:=strconv.Atoi(client.roomId)
			err=storage.CreateMessage(&types.Message{
				Content: msg.Content,
				Type:    msg.Type,
				UserID:  client.user.ID,
				ChatRoomId: uint(ChatRoomId),
			})
			if err!=nil{
				slog.Error("Error while creating message:", "error", err)
				continue
			}
			slog.Info("Message created successfully in db ")
        }
    }
}




// RabbitMQ


var MQChannel *amqp.Channel

func Init() error {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        slog.Error(fmt.Sprintf("Failed to connect to RabbitMQ: %s", err.Error()))
        return err // Return error to caller
    }

    ch, err := conn.Channel()
    if err != nil {
        slog.Error(fmt.Sprintf("Failed to create channel: %s", err.Error()))
        return err // Return error to caller
    }

    // Check if the channel is nil
    if ch == nil {
        slog.Error("RabbitMQ channel is nil")
        return fmt.Errorf("RabbitMQ channel is nil")
    }

    _, err = ch.QueueDeclare("chat", false, false, false, false, nil)
    if err != nil {
        slog.Error(fmt.Sprintf("Failed to declare queue: %s", err.Error()))
        return err // Return error to caller
    }

    MQChannel = ch // Set the global channel variable
    slog.Info("Successfully connected to RabbitMQ")
    return nil
}

func PublishQueue(body types.Message) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to marshal message %s ",err.Error()))
	}
	err=MQChannel.Publish("", "chat", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        bodyBytes,
	})
	if err!=nil{
		slog.Error(fmt.Sprintf("Failed to publish message %s ",err.Error()))
	}
	return nil
}

func ConsumeQueue(storage storage.Storage)  {
	messages, err :=MQChannel.Consume("chat", "", false, false, false, false, nil)
	if err!=nil{
		slog.Error(fmt.Sprintf("Failed to Consume queue %s ",err.Error()))
	}
	go func() {
		for d :=range messages{
			var msg types.Message
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				slog.Error("Failed to unmarshal message")
				continue
			}
			ClientMux.Lock()
			clients := rooms[strconv.Itoa(int(msg.ChatRoomId))]
			ClientMux.Unlock()

			for client := range  clients {
				select {
				case client.Send <-msg:
				slog.Info("Message sent successfully")
				default:
					slog.Warn("client's send channel is full, skipping message")
				}
			}
		}
	}()
}
