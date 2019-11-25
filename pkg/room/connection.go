package room

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

// Connection represents a websocket connection
type Connection struct {
	ID       string
	Conn     *websocket.Conn
	Nickname string
}

// NewConnection creates a new Connection
func NewConnection(c *websocket.Conn) *Connection {
	id := xid.New().String()

	return &Connection{
		id,
		c,
		id,
	}
}

// HandleConn handles the websocket connection
func (r *Room) HandleConn(conn *Connection) {
	r.AddConn(conn)
	defer r.RemoveConn(conn.ID)
	c := conn.Conn
	id := conn.ID

	c.WriteJSON(packet{identifyOp, packetIdentify{id, conn.Nickname}})
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("read(%s : %s): %s\n", id, conn.Nickname, err)
			return
		}
		var pckt incPacket
		json.Unmarshal(message, &pckt)
		switch pckt.Op {
		case heartBeatOp:
			log.Printf("heartbeat from %s", id)
		case claimOwnerShipOp:
			var msg map[string]string
			json.Unmarshal(message, &msg)
			fmt.Println(msg)
			if r.Secret == msg["secret"] {
				r.SetOwner(conn.ID)
			}
		case createOfferOp:
			if r.Owner == "" {
				continue
			}
			var offer packetCreateOffer
			json.Unmarshal(message, &offer)
			owner, err := r.GetConn(r.Owner)

			if err != nil {
				continue
			}
			owner.Conn.WriteJSON(packet{createOfferOp, packetCreateOffer{offer.Offer, conn.ID}})
		case sendAnswerOp:
			if conn.ID != r.Owner {
				continue
			}
			var answer packetSendAnswer
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
			owner.Conn.WriteJSON(packet{readyOp, packetReady{Recipient: conn.ID}})
		case iceCandidateOp:
			if r.Owner == "" {
				continue
			}
			var candidate packetIceCandidate
			json.Unmarshal(message, &candidate)
			if r.Owner == conn.ID {
				recipient, err := r.GetConn(candidate.Recipient)

				if err != nil {
					continue
				}

				recipient.Conn.WriteJSON(packet{iceCandidateOp, packetIceCandidate{r.Owner, candidate.Candidate}})
			} else {
				owner, err := r.GetConn(r.Owner)

				if err != nil {
					continue
				}

				owner.Conn.WriteJSON(packet{iceCandidateOp, packetIceCandidate{conn.ID, candidate.Candidate}})
			}
		}
	}
}
