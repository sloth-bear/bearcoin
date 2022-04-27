package blockchain

import (
	"sync"

	"github.com/sloth-bear/bearcoin/db"
	"github.com/sloth-bear/bearcoin/utils"
)

const (
	defaultDifficulty   int = 2
	difficultyInterval  int = 5
	blockIntervalMinute int = 2
	allowedRange        int = 2
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

var b *blockchain
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty

	persistBlockchain(b)
}

func persistBlockchain(b *blockchain) {
	db.SaveState(utils.ToBytes(b))
}

func Blocks(b *blockchain) []*Block {
	var blocks []*Block

	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)

		if block.PrevHash == "" {
			break
		}

		hashCursor = block.PrevHash
	}

	return blocks
}

func recalculateDifficulty(b *blockchain) int {
	allBlocks := Blocks(b)
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]

	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockIntervalMinute

	if actualTime < (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime > (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}

	return b.CurrentDifficulty
}

func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}

func FindTx(b *blockchain, txID string) *Tx {
	for _, tx := range Txs(b) {
		if tx.ID == txID {
			return tx
		}
	}
	return nil
}

func UTxOutsByAddress(b *blockchain, address string) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)

	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE" {
					break
				}

				if FindTx(b, input.TxID).TxOuts[input.Index].Address == address {
					creatorTxs[input.TxID] = true
				}
			}

			for index, output := range tx.TxOuts {
				if output.Address == address {
					if _, isSpent := creatorTxs[tx.ID]; !isSpent {
						uTxOut := &UTxOut{tx.ID, index, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}

	return uTxOuts
}

func BalanceByAddress(b *blockchain, address string) int {
	ownedTxOuts := UTxOutsByAddress(b, address)
	var amount int
	for _, txOut := range ownedTxOuts {
		amount += txOut.Amount
	}
	return amount
}

func Blockchain() *blockchain {
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		}
		state := db.State()

		if state == nil {
			b.AddBlock()
		} else {
			b.restore(state)
		}
	})

	return b
}
