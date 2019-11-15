package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/tortlewortle/egghtogether/internal/routes"
	"github.com/tortlewortle/egghtogether/pkg/connection"
	"github.com/tortlewortle/egghtogether/pkg/room"
)

var manager = room.NewManager()

var upgrader = websocket.Upgrader{} // use default options
func newRoom(w http.ResponseWriter, r *http.Request) {
	newRoom := room.New()
	manager.AddRoom(newRoom)

	out, _ := json.Marshal(newRoom)
	w.Write(out)
}

func joinRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["id"]
	currentRoom, err := manager.GetRoom(roomID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Room not found"))
		return
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	conn := connection.New(c)

	currentRoom.HandleConn(conn)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/TortleWortle/egghtogether/releases", 307)
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	r := mux.NewRouter()

	r.HandleFunc("/", redirect)
	r.HandleFunc("/newroom", newRoom)
	r.HandleFunc("/rooms/{id}", joinRoom)
	r.HandleFunc("/watch/{id}", routes.WatchRoute)
	// r.HandleFunc("/info", room.DebugInfoRoute)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
