package connection

import (
	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

// Connection represents a websocket connection
type Connection struct {
	ID   string
	Conn *websocket.Conn
}

// New creates a new Connection
func New(c *websocket.Conn) *Connection {
	id := xid.New().String()

	return &Connection{
		id,
		c,
	}
}
