package sse

import (
	"log"
	"sync"
)

type Hub struct {
	clients   map[string]*Client
	Broadcast chan string
	AddClient chan *Client
	mut       sync.RWMutex
}

// this function is to be called as a goroutine which continously sets up this as a service
func (h *Hub) Run() {
	for {
		select {
		case msg := <-h.Broadcast:
			log.Println("sending broadcast")
			h.mut.RLock()
			for _, c := range h.clients {
				c.MessageChan <- msg
			}
			h.mut.RUnlock()
		case client := <-h.AddClient:
			h.mut.RLock()
			h.clients[client.Id] = client
			h.mut.RUnlock()
			client.MessageChan <- "welcome"
			log.Println("client added")
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		clients:   map[string]*Client{},
		Broadcast: make(chan string),
		AddClient: make(chan *Client),
	}
}
