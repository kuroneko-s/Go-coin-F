package p2p

import (
	"fmt"
	"net/http"
	"time"

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

	conn, err := upgrader.Upgrade(w, r, nil)
	utils.HandleErr(err)
	fmt.Println(r.RemoteAddr)

	initPeer(conn, ip, openPort)
	time.Sleep(10 * time.Second)
	conn.WriteMessage(websocket.TextMessage, []byte("Hello from port 3000"))
}

// 요청하는 Peer가 사용
func AddPeer(address, port, openPort string) {
	// goLang에서 webSocket을 연결할땐 dialer가 필요하다. (코드상)
	// Header에다가 Authenticate tocken 같은거 넣어서 Upgrade시 인증받을 수 있게 진행할 수 있음
	// from :4000 -> :3000
	fmt.Printf("ws://%s:%s/ws\n", address, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleErr(err)
	initPeer(conn, address, port)
	time.Sleep(10 * time.Second)
	conn.WriteMessage(websocket.TextMessage, []byte("Hello from port 4000"))
}
