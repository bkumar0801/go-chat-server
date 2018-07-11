package main

import (
	"log"
	"net/http"

	"github.com/go-chat-server/context"
	"github.com/go-chat-server/handlers"
)

func main() {
	appCtx := context.NewAppContext()
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	http.HandleFunc("/echo", handlers.HandleConnections(appCtx))

	go handlers.HandleMessages(appCtx)

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
