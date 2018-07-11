package context

import (
	"net/http"

	"github.com/gorilla/websocket"
)

/*
Message ...
*/
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

/*
AppContenxt ...
*/
type AppContext struct {
	Users     map[string]int
	Clients   map[*websocket.Conn]bool
	Broadcast chan Message
	Upgrader  websocket.Upgrader
}

/*
NewAppContext ...
*/
func NewAppContext() *AppContext {
	return &AppContext{
		Users:     make(map[string]int),
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan Message),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}
