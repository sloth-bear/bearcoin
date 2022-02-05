package db

import (
	"github.com/boltdb/bolt"
	"github.com/sloth-bear/bearcoin/utils"
)

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"

	state = "state"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		utils.HandleErr(err)

		db = dbPointer

		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)

			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleErr(err)
	}

	return db
}

func SaveBlock(hash string, block []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), block)

		return err
	})

	utils.HandleErr(err)
}

func Block(hash string) []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})

	return data
}

// func FindAllBlocks() {
// 	DB().View(func(t *bolt.Tx) error {
// 		bucket := t.Bucket([]byte(blocksBucket))
// 		err := bucket.
// 	})
// }

func SaveBlockchain(blockchain []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(state), blockchain)

		return err
	})

	utils.HandleErr(err)
}

func BlockchainState() []byte {
	var data []byte

	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(state))
		return nil
	})

	return data
}
