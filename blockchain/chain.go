package blockchain

import (
	"sync"

	"github.com/sloth-bear/bearcoin/db"
	"github.com/sloth-bear/bearcoin/utils"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height

	b.persist()
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			state := db.BlockchainState()

			if state == nil {
				b.AddBlock("Genesis Block")
			} else {
				b.restore(state)
			}
		})
	}

	return b
}
