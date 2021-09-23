package p2p

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var Peers map[string]*peer = make(map[string]*peer)

type peer struct {
	conn  *websocket.Conn
	inbox chan []byte
}

func (p *peer) read() {
	// delete peer in case of error
	for {
		_, m, err := p.conn.ReadMessage() // block
		if err != nil {
			break
		}
		fmt.Printf("%s", m)
	}
}

func (p *peer) write() {
	for {
		m := <-p.inbox // block
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	p := &peer{
		conn: conn,
	}
	key := fmt.Sprintf("%s:%s", address, port)
	Peers[key] = p
	go p.read()

	return p
}
