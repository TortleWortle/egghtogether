package main

import (
	"encoding/json"
	"flag"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	assets "github.com/tortlewortle/egghtogether/internal/bindata"
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

	conn := room.NewConnection(c)

	currentRoom.HandleConn(conn)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/TortleWortle/egghtogether/releases", 307)
}

func serveAsset(w http.ResponseWriter, r *http.Request) {
	asset, err := assets.Asset(r.URL.Path[1:])
	if err != nil {
		asset, err = assets.Asset("index.html")
		if err != nil {
			panic("index.html not found D:")
		}
	}
	splat := strings.SplitN(r.URL.Path[1:], ".", -1)

	if len(splat) > 1 {
		w.Header().Set("Content-Type", mime.TypeByExtension("."+splat[1]))
	}

	w.Write(asset)
	// http.Redirect(w, r, "https://github.com/TortleWortle/egghtogether/releases", 307)
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	r := mux.NewRouter()

	r.HandleFunc("/api/newroom", newRoom)
	r.HandleFunc("/api/room/{id}/ws", joinRoom)
	r.HandleFunc("/api/room/{id}/info", manager.InfoRoute)
	r.NotFoundHandler = http.HandlerFunc(serveAsset)

	if os.Getenv("DEBUG") == "true" {
		r.HandleFunc("/api/info", manager.DebugInfoRoute)
	}

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	mux := &http.ServeMux{}
	mux.Handle("/metrics", promhttp.Handler())
	metricsSrv := &http.Server{
		Handler: mux,
		Addr:    "127.0.0.1:9000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go recordMetrics()
	go metricsSrv.ListenAndServe()

	log.Fatal(srv.ListenAndServe())
}
