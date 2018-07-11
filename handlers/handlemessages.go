package handlers

import (
	"log"

	"github.com/go-chat-server/context"
)

/*
HandleMessages ...
*/
func HandleMessages(ctx *context.AppContext) {
	for {
		msg := <-ctx.Broadcast
		if msg.Username != "Server" {
			log.Printf("%s: %s ", msg.Username, msg.Message)
		}
		for client := range ctx.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(ctx.Clients, client)
				log.Printf("error: %v", err)
			}
		}
	}
}
