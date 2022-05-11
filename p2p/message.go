package p2p

import (
	"encoding/json"

	"github.com/sloth-bear/bearcoin/blockchain"
	"github.com/sloth-bear/bearcoin/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func (m *Message) setPayload(p interface{}) {
	pAsBytes, err := json.Marshal(p)
	utils.HandleErr(err)

	m.Payload = pAsBytes
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{Kind: kind}
	m.setPayload(payload)

	mAsJson, err := json.Marshal(m)
	utils.HandleErr(err)

	return mAsJson
}

func sendNewestBlock(p *peer) {
	b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)

	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m
}
