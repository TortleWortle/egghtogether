package room

import (
	"encoding/json"
	"log"

	"github.com/tortlewortle/egghtogether/pkg/connection"
)

// HandleConn handles the websocket connection
func (r *Room) HandleConn(conn *connection.Connection) {
	r.AddConn(conn)
	defer r.RemoveConn(conn.ID)
	c := conn.Conn
	id := conn.ID

	c.WriteJSON(packet{identifyOp, identifyPacket{id}})
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("read(%s): %s\n", id, err)
			return
		}
		var msg map[string]interface{}
		json.Unmarshal(message, &msg)
		log.Printf("recv(%s): %s", id, msg["op"])
		switch msg["op"] {
		case heartBeatOp:
			log.Printf("heartbeat from %s", id)
		case claimOwnerShipOp:
			if r.Secret == msg["secret"] {
				r.SetOwner(conn.ID)
			}
		case createOfferOp:
			if r.Owner == "" {
				continue
			}
			var offer createOfferPacket
			json.Unmarshal(message, &offer)
			owner, err := r.GetConn(r.Owner)

			if err != nil {
				continue
			}
			owner.Conn.WriteJSON(packet{createOfferOp, createOfferPacket{offer.Offer, conn.ID}})
		case sendAnswerOp:
			if conn.ID != r.Owner {
				continue
			}
			var answer sendAnswerPacket
			json.Unmarshal(message, &answer)
			recipient, err := r.GetConn(answer.Recipient)

			if err != nil {
				continue
			}

			recipient.Conn.WriteJSON(packet{sendAnswerOp, answer})
		case readyOp:
			if r.Owner == "" {
				continue
			}
			owner, err := r.GetConn(r.Owner)

			if err != nil {
				continue
			}
			owner.Conn.WriteJSON(packet{readyOp, readyPacket{Recipient: conn.ID}})
		case iceCandidateOp:
			if r.Owner == "" {
				continue
			}
			var candidate iceCandidatePacket
			json.Unmarshal(message, &candidate)
			if r.Owner == conn.ID {
				recipient, err := r.GetConn(candidate.Recipient)

				if err != nil {
					continue
				}

				recipient.Conn.WriteJSON(packet{iceCandidateOp, iceCandidatePacket{r.Owner, candidate.Candidate}})
			} else {
				owner, err := r.GetConn(r.Owner)

				if err != nil {
					continue
				}

				owner.Conn.WriteJSON(packet{iceCandidateOp, iceCandidatePacket{conn.ID, candidate.Candidate}})
			}
		}
	}
}
