package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/goLangCoin/blockchain"
	"github.com/goLangCoin/utils"
)

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
)

type MessageKind int

type Message struct {
	Kind    MessageKind
	Payload []byte // どのタイプでも出来る。null or []Block or block
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	// interfaceのpayloadをどのタイプでも変換してくれるfunc
	m := Message{
		Kind:    kind,
		Payload: utils.ToJson(payload),
	}

	return utils.ToJson(m)
}

func sendNewestBlock(p *peer) {
	b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)

	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m // inbox는 값을 보낼 수 있는 박스임
}

func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksResponse, blockchain.Blocks(blockchain.Blockchain()))
	p.inbox <- m
}

// response에서 동작
func handleMsg(m *Message, p *peer) {
	// fmt.Printf("Peers: %s, sent a meesage %d", p.key, m.Payload)
	switch m.Kind {
	case MessageNewestBlock:
		var payload blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
		utils.HandleErr(err)
		if payload.Height >= b.Height {
			requestAllBlocks(p)
		} else {
			// request보다 최신의 block을 가지고 있다면 상대방한테 최신 블럭에 대한 정보를 보냄
			sendNewestBlock(p)
		}
	case MessageAllBlocksRequest:
		sendAllBlocks(p)
	case MessageAllBlocksResponse:
		var payload *[]blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		fmt.Println(payload)
	}

}
