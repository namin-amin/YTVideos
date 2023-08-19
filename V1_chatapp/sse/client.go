package sse

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Client struct {
	Id          string
	writter     http.ResponseWriter
	MessageChan chan string
}

func (c *Client) RunSSE() {
	c.writter.Header().Add("Content-Type", "text/event-stream")
	c.writter.Header().Add("Cache-Control", "no-cache")
	c.writter.Header().Add("Connection", "keep-alive")
	c.writter.Header().Add("Transfer-Encoding", "chunked")

	flusher := c.writter.(http.Flusher)

	//intitial connection setup
	_, err := fmt.Fprintf(c.writter, "data:%s\n\n", "Hello welcome")
	if err != nil {
		log.Println("something went wrong while connecting to client")
		return
	}

	flusher.Flush()

	//server sent event connection
	//dont close connection but continously listen to the message channel to send it
	for {
		select {
		case msg := <-c.MessageChan:
			_, err := fmt.Fprintf(c.writter, "data:%s\n\n", msg)
			if err != nil {
				log.Println("something went wrong while connecting to client")
				break
			}
			flusher.Flush()
			//Todo add other cases this if only one case better to use forloop range
		}
	}

}

func NewClient(w http.ResponseWriter) *Client {
	return &Client{
		writter:     w,
		Id:          uuid.NewString(),
		MessageChan: make(chan string),
	}
}
