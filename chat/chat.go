/* */

package chat

import (
	"casServer/login"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var cLog *log.Logger

// Message object
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// IngameChat ...
type IngameChat struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
	upgrader  websocket.Upgrader
	cLog      *log.Logger
}

// NewChat ...
func NewChat(chatLogFile string) *IngameChat {
	createCustomChatLogFile(chatLogFile)
	return &IngameChat{
		make(map[*websocket.Conn]bool),
		make(chan Message),
		websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		cLog,
	}
}

// HandleChat ...
func (ch *IngameChat) HandleChat(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := ch.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ERROR 1 Chat => ", err)
	}
	defer ws.Close()

	// Register our new client
	ch.clients[ws] = true

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error 2 Chat=> %v", err)
			delete(ch.clients, ws)
			break
		}
		// check anonymous user name is valid
		r.ParseForm()
		cookie, _ := r.Cookie(login.CookieName)
		msg.Username = strings.Split(cookie.Value, ":")[0]
		// send message to broadcast channel
		ch.cLog.Println(msg.Username, " -> ", msg.Message)
		ch.broadcast <- msg
	}
}

// BroadcastMessages ...
func (ch *IngameChat) BroadcastMessages() {
	for {
		// Grab message from broadcast channel
		msg := <-ch.broadcast
		// Send it out to every client connected
		for client := range ch.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(ch.clients, client)
			}
		}
	}
}
