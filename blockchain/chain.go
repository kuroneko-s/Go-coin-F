package blockchain

import (
	"fmt"
	"sync"

	"github.com/goLangCoin/db"
	"github.com/goLangCoin/utils"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5 // 생성주기 조건 5개마다 재조정
	blockInterval      int = 2 // 생성주기 2분
	allowedRange       int = 2
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentdifficulty"`
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
	block := createBlock(data, b.NewestHash, b.Height)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

// 난이도 재조정
func (b *blockchain) recalculateDifficulty() int {
	// 5개 생성하는데에 10분 걸려야함.
	// 기대값 - 2분에 한개 생성, 5개 생성시 주기 계산해서 난이도 조정
	// 가장 최근 블럭 가져오고
	// 이전 블럭 가져옴 ( 5번전꺼 기준이니깐)
	// 최근 - (최근-5)에 대한 값을 구해서 시간을 구한다.
	allBlocks := b.Blocks()
	newestBlock := allBlocks[0]                                                         // 가장 최근 블럭
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]                            // 가장 최근 재설정된 블럭
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60) // 5개를 생성하는데 걸린 시간 (분)
	expectedtime := difficultyInterval * blockInterval                                  // 5개를 생성하는데 걸려야 하는 시간
	if actualTime <= (expectedtime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedtime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func (b *blockchain) difficulty() int {
	// bicoin checks in created 2016
	if b.Height == 0 {
		return defaultDifficulty // default
	} else if b.Height%difficultyInterval == 0 {
		// 5개마다 체크해서 난이도 재조정
		return b.recalculateDifficulty()
	}
	// 이전 블럭의 difficulty를 그대로 가져옴
	return b.CurrentDifficulty
}

// Singleton pattern
func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			// 맨처음엔 마지막 블럭이 어떤값을 이루고 있는지를 모르고잇으니깐
			b = &blockchain{
				Height: 0,
			}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis Block")
			} else {
				fmt.Println("Restoring...")
				b.restore(checkpoint)
			}
		})
	}
	return b
}
