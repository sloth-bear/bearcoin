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
	fmt.Printf("Sending newest block to %s\n", p.key)
	b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)

	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	fmt.Printf("Sending all blocks to %s\n", p.key)
	m := makeMessage(MessageAllBlocksResponse, blockchain.Blocks(blockchain.Blockchain()))
	p.inbox <- m
}

func requestAllBlocks(p *peer) {
	fmt.Printf("Requesting all blocks to %s\n", p.key)
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

func handleMsg(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		fmt.Printf("Received the newest block from %s\n", p.key)
		var payload blockchain.Block
		err := json.Unmarshal(m.Payload, &payload)
		utils.HandleErr(err)

		b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
		utils.HandleErr(err)

		if payload.Height >= b.Height {
			requestAllBlocks(p)
		} else {
			sendNewestBlock(p)
		}
	case MessageAllBlocksRequest:
		fmt.Printf("%s wants all the blocks\n", p.key)
		sendAllBlocks(p)
	case MessageAllBlocksResponse:
		fmt.Printf("Received all the blocks from %s\n", p.key)
		var payload []*blockchain.Block

		err := json.Unmarshal(m.Payload, &payload)
		utils.HandleErr(err)
	}
}
