package socket

import (
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Conn      *websocket.Conn
	timestamp time.Time
}

func NewConnection(conn *websocket.Conn) *Connection {
	return &Connection{
		Conn:      conn,
		timestamp: time.Now(),
	}
}
