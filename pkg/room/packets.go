package room

type packet struct {
	Op   string      `json:"op"`
	Data interface{} `json:"data"`
}
type identifyPacket struct {
	ID string `json:"id"`
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
