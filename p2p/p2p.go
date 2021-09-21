package p2p

import (
	"fmt"
	"net/http"

	"github.com/goLangCoin/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var conns []*websocket.Conn

func Upgrade(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	utils.HandleErr(err)

	conns = append(conns, conn)
	fmt.Println("Connections - ", conns)
	for {
		_, p, err := conn.ReadMessage() // 여기서 blocking 해줌
		if err != nil {
			conn.Close()
			break
		}
		for _, aConn := range conns {
			if aConn != conn {
				utils.HandleErr(aConn.WriteMessage(websocket.TextMessage, p))
			}
		}
	}
}
