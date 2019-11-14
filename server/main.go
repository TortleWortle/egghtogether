package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

var upgrader = websocket.Upgrader{} // use default options
func newRoom(w http.ResponseWriter, r *http.Request) {
	room := NewRoom()
	addRoom(room)

	out, _ := json.Marshal(room)
	w.Write(out)
}

func joinRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["id"]
	room, err := getRoom(roomID)

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
	id := xid.New().String()

	conn := &Connection{
		ID:   id,
		conn: c,
	}

	room.addConn(conn)
	defer room.removeConn(id)

	handleConn(room, conn)
}

const (
	identifyOp       = "IDENTIFY"
	messageOp        = "MESSAGE"
	claimOwnerShipOp = "CLAIM_OWNERSHIP"
	createOfferOp    = "CREATE_OFFER"
	sendAnswerOp     = "SEND_ANSWER"
	readyOp          = "READY"
	iceCandidateOp   = "ICE_CANDIDATE"
)

type packet struct {
	Op   string      `json:"op"`
	Data interface{} `json:"data"`
}
type identifyPacket struct {
	ID string `json:"id"`
}
type messagePacket struct {
	Msg string `json:"msg"`
}
type createOfferPacket struct {
	Offer     string `json:"offer"`
	Recipient string `json:"recipient"`
}
type sendAnswerPacket struct {
	Answer    string `json:"answer"`
	Recipient string `json:"recipient"`
}
type readyPacket struct {
	Recipient string `json:"recipient"`
}
type iceCandidatePacket struct {
	Recipient string `json:"recipient"`
	Candidate string `json:"candidate"`
}

func handleConn(room *Room, conn *Connection) {
	c := conn.conn
	id := conn.ID
	c.WriteJSON(packet{identifyOp, identifyPacket{id}})
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("read(%s): %s\n", id, err)
			return
		}
		log.Printf("recv(%s): %s", id, message)
		var msg map[string]interface{}
		json.Unmarshal(message, &msg)
		switch msg["op"] {
		case claimOwnerShipOp:
			if room.Secret == msg["secret"] {
				room.setOwner(conn.ID)
			}
		case createOfferOp:
			if room.Owner == "" {
				continue
			}
			var offer createOfferPacket
			json.Unmarshal(message, &offer)
			owner, err := room.getConn(room.Owner)

			if err != nil {
				continue
			}
			owner.conn.WriteJSON(packet{createOfferOp, createOfferPacket{offer.Offer, conn.ID}})
		case sendAnswerOp:
			if conn.ID != room.Owner {
				continue
			}
			var answer sendAnswerPacket
			json.Unmarshal(message, &answer)
			recipient, err := room.getConn(answer.Recipient)

			if err != nil {
				continue
			}

			recipient.conn.WriteJSON(packet{sendAnswerOp, answer})
		case readyOp:
			if room.Owner == "" {
				continue
			}
			owner, err := room.getConn(room.Owner)

			if err != nil {
				continue
			}
			owner.conn.WriteJSON(packet{readyOp, readyPacket{Recipient: conn.ID}})
		case iceCandidateOp:
			if room.Owner == "" {
				continue
			}
			var candidate iceCandidatePacket
			json.Unmarshal(message, &candidate)
			if room.Owner == conn.ID {
				recipient, err := room.getConn(candidate.Recipient)

				if err != nil {
					continue
				}

				recipient.conn.WriteJSON(packet{iceCandidateOp, iceCandidatePacket{room.Owner, candidate.Candidate}})
			} else {
				owner, err := room.getConn(room.Owner)

				if err != nil {
					continue
				}

				owner.conn.WriteJSON(packet{iceCandidateOp, iceCandidatePacket{conn.ID, candidate.Candidate}})
			}
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	r := mux.NewRouter()

	r.HandleFunc("/newroom", newRoom)
	r.HandleFunc("/rooms/{id}", joinRoom)
	r.HandleFunc("/watch/{id}", watchRoute)
	r.HandleFunc("/info", roomInfo)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
