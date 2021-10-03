package p2p

import (
	"fmt"
	"net/http"

	"github.com/goLangCoin/blockchain"
	"github.com/goLangCoin/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// 응답하는 Peer가 사용
func Upgrade(w http.ResponseWriter, r *http.Request) {
	// from :3000 -> :4000 응답(upgrade)
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	openPort := r.URL.Query().Get("openPort")

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}

	fmt.Printf("%s wants an upgrade\n", openPort)
	conn, err := upgrader.Upgrade(w, r, nil)
	utils.HandleErr(err)

	initPeer(conn, ip, openPort)
}

// 요청하는 Peer가 사용
func AddPeer(address, port, openPort string, broadcast bool) {
	// goLang에서 webSocket을 연결할땐 dialer가 필요하다. (코드상)
	// Header에다가 Authenticate tocken 같은거 넣어서 Upgrade시 인증받을 수 있게 진행할 수 있음
	// from :4000 -> :3000
	fmt.Printf("%s want to connect to port %s\n", openPort, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleErr(err)
	p := initPeer(conn, address, port)
	if broadcast {
		broadcastNewPeer(p)
		return
	}
	sendNewestBlock(p) // 연결에 성공했을 경우 Three Hand Shake처럼 가장 최신의 Block에 대한 값을 보냄
}

func BroadcaseNewBlock(b *blockchain.Block) {
	Peers.m.Lock()
	defer Peers.m.Unlock()

	for _, p := range Peers.v {
		notifyNewBlock(b, p)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	Peers.m.Lock()
	defer Peers.m.Unlock()

	for _, p := range Peers.v {
		notifyNewTx(tx, p)
	}
}

func broadcastNewPeer(newPeer *peer) {
	Peers.m.Lock()
	defer Peers.m.Unlock()

	for k, p := range Peers.v {
		if k != newPeer.key {
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(payload, p)
		}
	}
}
