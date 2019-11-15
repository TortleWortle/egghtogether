package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

var rooms = make(map[string]*Room)
var roomlock sync.Mutex

func getRoom(id string) (*Room, error) {
	roomlock.Lock()
	defer roomlock.Unlock()

	room, ok := rooms[id]

	if !ok {
		return &Room{}, errors.New("Room does not exist")
	}
	return room, nil
}

func addRoom(room *Room) {
	roomlock.Lock()
	defer roomlock.Unlock()
	rooms[room.ID] = room
}

func removeRoomByID(id string) {
	roomlock.Lock()
	defer roomlock.Unlock()
	delete(rooms, id)
}

func removeRoom(room *Room) {
	roomlock.Lock()
	defer roomlock.Unlock()
	delete(rooms, room.ID)
}

// Connection represents a websocket connection
type Connection struct {
	ID   string
	conn *websocket.Conn
}

// Room represents a room
type Room struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	Owner       string `json:"owner"`
	connections map[string]*Connection
	connlock    sync.Mutex
}

// NewRoom makes a new room
func NewRoom() *Room {
	id := xid.New().String()
	secret, err := GenerateRandomString(64)

	if err != nil {
		panic(err)
	}
	room := &Room{
		ID:          id,
		connections: make(map[string]*Connection),
		Secret:      secret,
	}
	return room
}

func (r *Room) getConn(id string) (*Connection, error) {
	r.connlock.Lock()
	defer r.connlock.Unlock()
	conn, ok := r.connections[id]

	if !ok {
		return &Connection{}, errors.New("Connection doesn't exist")
	}
	return conn, nil
}

func (r *Room) setOwner(id string) {
	r.connlock.Lock()
	defer r.connlock.Unlock()
	r.Owner = id
}

func (r *Room) removeConn(id string) {
	r.connlock.Lock()
	defer r.connlock.Unlock()
	delete(r.connections, id)
	if len(r.connections) == 0 || id == r.Owner {
		removeRoom(r)
	}
}

func (r *Room) addConn(c *Connection) {
	r.connlock.Lock()
	defer r.connlock.Unlock()
	r.connections[c.ID] = c
}

type roomInfoStruct struct {
	ID              string `json:"id"`
	ConnectionCount int    `json:"connection_count"`
	Owner           string `json:"owner"`
}

func roomInfo(w http.ResponseWriter, req *http.Request) {
	infos := make([]roomInfoStruct, 0)

	for _, room := range rooms {
		infos = append(infos, roomInfoStruct{
			ID:              room.ID,
			ConnectionCount: len(room.connections),
			Owner:           room.Owner,
		})
	}

	out, err := json.Marshal(infos)
	if err != nil {
		w.Write([]byte("Err"))
		return
	}
	w.Write(out)
}
