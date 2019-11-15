package room

import (
	"errors"
	"sync"

	"github.com/rs/xid"
	"github.com/tortlewortle/egghtogether/internal/util"
	"github.com/tortlewortle/egghtogether/pkg/connection"
)

// Room represents a room
type Room struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	Owner       string `json:"owner"`
	connections map[string]*connection.Connection
	connlock    sync.Mutex
}

// New makes a new room
func New() *Room {
	id := xid.New().String()
	secret, err := util.GenerateRandomString(64)

	if err != nil {
		panic(err)
	}
	room := &Room{
		ID:          id,
		connections: make(map[string]*connection.Connection),
		Secret:      secret,
	}
	return room
}

// GetConn gets a connection from a room
func (r *Room) GetConn(id string) (*connection.Connection, error) {
	r.connlock.Lock()
	defer r.connlock.Unlock()
	conn, ok := r.connections[id]

	if !ok {
		return &connection.Connection{}, errors.New("Connection doesn't exist")
	}
	return conn, nil
}

// SetOwner sets the owner of the room
func (r *Room) SetOwner(id string) {
	r.connlock.Lock()
	defer r.connlock.Unlock()
	r.Owner = id
}

// RemoveConn removes a connection from a room
func (r *Room) RemoveConn(id string) {
	r.connlock.Lock()
	defer r.connlock.Unlock()
	delete(r.connections, id)
}

// AddConn adds a connection to a room
func (r *Room) AddConn(c *connection.Connection) {
	r.connlock.Lock()
	defer r.connlock.Unlock()
	r.connections[c.ID] = c
}
