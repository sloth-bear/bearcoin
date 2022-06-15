package blockchain

import (
	"reflect"
	"testing"

	"github.com/sloth-bear/bearcoin/utils"
)

func TestCreateBlock(t *testing.T) {
	dbStorage = fakeDB{}
	Mempool().Txs["test"] = &Tx{}

	b := createBlock("TEST_HASH", 1, 1)
	if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
		t.Error("createBlock() should return an instance of a block")
	}
}

func TestFindBlock(t *testing.T) {
	t.Run("Should return error if block doesn't exist", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func(hash string) []byte {
				return nil
			},
		}

		_, err := FindBlock("TEST_HASH")
		if err == nil {
			t.Error("The error should be return")
		}
	})

	t.Run("Should return block not found error", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func(hash string) []byte {
				b := &Block{
					Height: 1,
				}
				return utils.ToBytes(b)
			},
		}

		block, _ := FindBlock("TEST_HASH")
		if block.Height != 1 {
			t.Error("The block should be found")
		}
	})
}
