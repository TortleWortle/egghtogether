package room

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Manager manages rooms
type Manager struct {
	rooms map[string]*Room
	lock  *sync.Mutex
}

// NewManager makes a new manager
func NewManager() Manager {
	return Manager{
		make(map[string]*Room),
		&sync.Mutex{},
	}
}

// GetRoom gets a room by id
func (m *Manager) GetRoom(id string) (*Room, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	room, ok := m.rooms[id]

	if !ok {
		return &Room{}, errors.New("Room does not exist")
	}
	return room, nil
}

// AddRoom adds a room to local store
func (m *Manager) AddRoom(room *Room) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.rooms[room.ID] = room
}

// RemoveRoomByID removes a room by ID
func (m *Manager) RemoveRoomByID(id string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.rooms, id)
}

// RemoveRoom removes a room
func (m *Manager) RemoveRoom(room *Room) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.rooms, room.ID)
}

type roomInfoStruct struct {
	ID              string `json:"id"`
	ConnectionCount int    `json:"connection_count"`
	Owner           string `json:"owner"`
}

// InfoRoute displays info for a room
func (m *Manager) InfoRoute(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	roomID := vars["id"]
	currentRoom, err := m.GetRoom(roomID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	info := roomInfoStruct{
		ID:              currentRoom.ID,
		ConnectionCount: len(currentRoom.connections),
		Owner:           currentRoom.Owner,
	}

	out, err := json.Marshal(info)
	w.Write(out)
}

// DebugInfoRoute is a temporary debugging route
func (m *Manager) DebugInfoRoute(w http.ResponseWriter, req *http.Request) {
	infos := make([]roomInfoStruct, 0)

	for _, room := range m.rooms {
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

// TotalRooms returns the total amount of rooms
func (m *Manager) TotalRooms() int {
	return len(m.rooms)
}

// TotalConnections returns the total amount of connections
func (m *Manager) TotalConnections() int {
	total := 0

	for id := range m.rooms {
		room, err := m.GetRoom(id)

		if err != nil {
			continue
		}
		total += len(room.connections)
	}

	return total
}
