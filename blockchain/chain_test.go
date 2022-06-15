package blockchain

import (
	"sync"
	"testing"

	"github.com/sloth-bear/bearcoin/utils"
)

type fakeDB struct {
	fakeLoadChain func() []byte
	fakeFindBlock func(hash string) []byte
}

func (f fakeDB) FindBlock(hash string) []byte {
	return f.fakeFindBlock(hash)
}

func (fakeDB) SaveBlock(hash string, block []byte) {}

func (fakeDB) DeleteAllBlocks() {}

func (fakeDB) SaveChain(blockchain []byte) {}

func (f fakeDB) LoadChain() []byte {
	return f.fakeLoadChain()
}

func TestBlockchain(t *testing.T) {
	t.Run("Should create blockchain if not exists", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeLoadChain: func() []byte {
				return nil
			},
		}

		bc := Blockchain()
		if bc.Height != 1 {
			t.Error("Blockchain() should create a blockchain")
		}
	})

	t.Run("Should restore blockchain if exists", func(t *testing.T) {
		once = *new(sync.Once)
		dbStorage = fakeDB{
			fakeLoadChain: func() []byte {
				bc := &blockchain{Height: 1, NewestHash: "NewestHash", CurrentDifficulty: 1}
				return utils.ToBytes(bc)
			},
		}

		bc := Blockchain()
		if bc.NewestHash != "NewestHash" {
			t.Errorf("Blockchain() should restore a blockchain with a newestHash of %s, god %s", "NewestHash", bc.NewestHash)
		}
	})
}
