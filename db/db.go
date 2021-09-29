package db

import (
	"fmt"
	"os"

	"github.com/goLangCoin/utils"
	bolt "go.etcd.io/bbolt"
)

// bolt = bucket && DB = table
var db *bolt.DB

// bucket > block
// bucket > data

const (
	dbName       = "blockchain"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

func getDbName() string {
	port := os.Args[2][7:]
	return fmt.Sprintf("%s_%s.db", dbName, port)
}

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(getDbName(), 0600, nil)
		utils.HandleErr(err)
		db = dbPointer
		// transaction
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			utils.HandleErr(err)
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

// hash : data
func SaveBlock(hash string, data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

// checkpoint : data
// 오직 마지막 블럭에 대한 정보만 가지고 있음
func SaveBlockchain(data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

func Checkpoint() []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
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

func Close() {
	DB().Close()
}

func EmptyBlocks() {
	// 블록들에 대한 정보를 업데이트 할 때 DB를 비워둬야 해서 싹다 지우는 코드
	DB().Update(func(t *bolt.Tx) error {
		utils.HandleErr(t.DeleteBucket([]byte(blocksBucket)))
		_, err := t.CreateBucket([]byte(blocksBucket))
		utils.HandleErr(err)
		return nil
	})
}
