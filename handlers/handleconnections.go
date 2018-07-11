package handlers

import (
	"log"
	"net/http"

	"github.com/go-chat-server/context"
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
HandleConnections ...
*/
func HandleConnections(ctx *context.AppContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user string
		var msg context.Message

		ws, err := ctx.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer ws.Close()

		ctx.Clients[ws] = true

		err = ws.ReadJSON(&msg)
		if err != nil {
			delete(ctx.Clients, ws)
		} else {
			user = msg.Username
			makeUserOnline(ctx, &msg, user)
		}

		for {
			ctx.Broadcast <- msg
			err := ws.ReadJSON(&msg)
			if err != nil {
				delete(ctx.Clients, ws)
				makeUserOffline(ctx, &msg, user)
				break
			}
		}
	}
}

func makeUserOnline(ctx *context.AppContext, msg *context.Message, user string) {
	if val, ok := ctx.Users[user]; ok {
		ctx.Users[user] = val + 1
		msg.Username = "Server"
		msg.Message = user + " is already logged in"
	} else {
		ctx.Users[user] = 1
		msg.Message = user + " is online"
		msg.Username = "Server"
		log.Printf("Server: %s is online", user)
	}
}

func makeUserOffline(ctx *context.AppContext, msg *context.Message, user string) {
	if ctx.Users[user] == 1 {
		delete(ctx.Users, user)
		msg.Message = user + " is offline"
		msg.Username = "Server"
		go func() {
			ctx.Broadcast <- *msg
		}()
		log.Printf("Server : %s is offline", user)
	} else {
		val := ctx.Users[user]
		ctx.Users[user] = val - 1
	}
}
