package client

import "golang.org/x/net/websocket"

type Connection struct {
	conn *websocket.Conn
}
