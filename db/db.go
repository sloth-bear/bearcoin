package db

import (
	"fmt"
	"os"

	"github.com/sloth-bear/bearcoin/utils"
	bolt "go.etcd.io/bbolt"
)

const (
	dbName       = "blockchain"
	dataBucket   = "data"
	blocksBucket = "blocks"

	state = "state"
)

var db *bolt.DB

type DB struct{}

func (DB) FindBlock(hash string) []byte {
	return findBlock(hash)
}
func (DB) SaveBlock(hash string, block []byte) {
	saveBlock(hash, utils.ToBytes(block))
}
func (DB) DeleteAllBlocks() {
	deleteAllBlocks()
}
func (DB) SaveChain(blockchain []byte) {
	saveChain(blockchain)
}
func (DB) LoadChain() []byte {
	return loadChain()
}

func getDbName() string {
	port := os.Args[2][6:]
	return fmt.Sprintf("%s_%s.db", dbName, port)
}

func InitDB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(getDbName(), 0600, nil)
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

func Close() {
	db.Close()
}

func saveBlock(hash string, block []byte) {
	err := db.Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), block)

		return err
	})

	utils.HandleErr(err)
}

func findBlock(hash string) []byte {
	var data []byte
	db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})

	return data
}

func saveChain(blockchain []byte) {
	err := db.Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(state), blockchain)

		return err
	})

	utils.HandleErr(err)
}

func loadChain() []byte {
	var data []byte

	db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(state))
		return nil
	})

	return data
}

func deleteAllBlocks() {
	db.Update(func(t *bolt.Tx) error {
		utils.HandleErr(t.DeleteBucket([]byte(blocksBucket)))
		_, err := t.CreateBucket([]byte(blocksBucket))
		utils.HandleErr(err)
		return nil
	})
}
