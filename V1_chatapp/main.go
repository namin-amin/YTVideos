package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/namin-amin/chatapp/sse"
)

func main() {
	router := chi.NewRouter()

	hub := sse.NewHub()
	go hub.Run()

	//test path
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("hello world"))
	})

	//sse part

	//This is the endpoint used to connect for getting SSE connection
	router.Get("/sse", func(w http.ResponseWriter, r *http.Request) {
		client := sse.NewClient(w)
		hub.AddClient <- client
		client.RunSSE()
	})

	//sending message
	//should use post but for demo iam showing with get
	router.Get("/msg", func(w http.ResponseWriter, r *http.Request) {
		hub.Broadcast <- time.Now().Local().String()
		w.WriteHeader(200)
		w.Write([]byte("sent message"))
	})

	http.ListenAndServe(":3000", router)
}

//Bye
