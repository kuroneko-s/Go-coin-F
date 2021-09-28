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

func sendNEwestBlock(p *peer) {
	b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)

	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m // inbox는 값을 보낼 수 있는 박스임
}

func handleMsg(m *Message, p *peer) {
	// fmt.Printf("Peers: %s, sent a meesage %d", p.key, m.Payload)
	switch m.Kind {
	case MessageNewestBlock:
		var payload blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		fmt.Println(payload)
	case MessageAllBlocksRequest:
	case MessageAllBlocksResponse:

	}

}
