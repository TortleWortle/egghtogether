package room

type incPacket struct {
	Op   string `json:"op"`
	Data string `json:"data"`
}

type packet struct {
	Op   string      `json:"op"`
	Data interface{} `json:"data"`
}

type packetIdentify struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
}

type packetCreateOffer struct {
	Offer     string `json:"offer"`
	Recipient string `json:"recipient"`
}

type packetSendAnswer struct {
	Answer    string `json:"answer"`
	Recipient string `json:"recipient"`
}

type packetReady struct {
	Recipient string `json:"recipient"`
}

type packetIceCandidate struct {
	Recipient string `json:"recipient"`
	Candidate string `json:"candidate"`
}

type packetSendMessage struct {
	Message string `json:"message"`
}

type packetChatMessage struct {
	Message  string `json:"message"`
	Sender   string `json:"sender"`
	Nickname string `json:"nickname"`
}

type packetSetNickname struct {
	Nickname string `json:"nickname"`
}

type packetReceiveNickname struct {
	Sender   string `json:"sender"`
	Nickname string `json:"nickname"`
	OldNick  string `json:"old"`
}
