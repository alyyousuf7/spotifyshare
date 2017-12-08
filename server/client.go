package server

import (
	"context"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	ctx    context.Context
	hub    *Hub
	conn   *websocket.Conn
	sendCh chan []byte
	readCh chan []byte
	name   string
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		ctx:    context.TODO(),
		hub:    hub,
		conn:   conn,
		sendCh: make(chan []byte),
		readCh: make(chan []byte),
	}
}

func (c *Client) SetName(name string) {
	c.name = name
}

func (c *Client) Name() string {
	return c.name
}

func (c *Client) ReadLoop() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Errorf("Read error: %s", err)
		}

		c.readCh <- message
	}
}

func (c *Client) WriteLoop() {
	for {
		message := <-c.sendCh
		err := c.conn.WriteMessage(websocket.TextMessage, message)

		if err != nil {
			log.Error(err)
		}
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}
