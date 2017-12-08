package server

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Hub struct {
	startOnce *sync.Once
	clients   map[*Client]bool

	register chan *Client
}

func NewHub(ctx context.Context) *Hub {
	return &Hub{
		startOnce: &sync.Once{},
		clients:   make(map[*Client]bool),
		register:  make(chan *Client),
	}
}

func (h *Hub) Run() {
	h.startOnce.Do(h.run)
}

func (h *Hub) run() {
	go func() {
		for {
			select {
			case client := <-h.register:
				log.Info("Client register request")
				h.clients[client] = true

				// do not attempt to read or write before these loops
				go client.ReadLoop()
				go client.WriteLoop()

				log.Infof("Client joined!")
				h.Broadcast([]byte("Client joined"))
			default:
				// Do nothing
			}
		}
	}()
}

func (h *Hub) Register(conn *websocket.Conn) {
	client := NewClient(h, conn)
	h.register <- client
}

func (h *Hub) Broadcast(message []byte) {
	for c := range h.clients {
		c.sendCh <- message
	}
}

func (h *Hub) Close() error {
	for client := range h.clients {
		if err := client.Close(); err != nil {
			return err
		}
	}

	return nil
}
