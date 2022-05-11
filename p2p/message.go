package p2p

import (
	"encoding/json"
	"fmt"

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

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.JsonToBytes(payload),
	}

	return utils.JsonToBytes(m)
}

func sendNewestBlock(p *peer) {
	b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)

	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m
}

func handleMsg(m *Message, from *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		var payload blockchain.Block
		err := json.Unmarshal(m.Payload, &payload)
		utils.HandleErr(err)
		fmt.Println(payload)
	}
}
