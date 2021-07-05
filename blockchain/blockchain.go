package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

type blockchain struct {
	blocks []*Block // copy 방지
}

var b *blockchain
var once sync.Once
var ErrLenHandle = errors.New("Not found Block")

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlocks := len(AllBlocks())
	if totalBlocks == 0 {
		return ""
	}
	return AllBlocks()[totalBlocks-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(), len(AllBlocks()) + 1}
	newBlock.calculateHash()

	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

// Singleton pattern
func GetBlockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}

func AllBlocks() []*Block {
	return GetBlockchain().blocks
}

func (b blockchain) GetBlock(id int) (*Block, error) {
	if id > len(b.blocks) || id == 0 {
		return nil, ErrLenHandle
	}
	return b.blocks[id-1], nil
}
