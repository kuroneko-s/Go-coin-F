package p2p

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var Peers map[string]*peer = make(map[string]*peer)

type peer struct {
	conn    *websocket.Conn
	inbox   chan []byte
	key     string
	address string
	port    string
}

func (p *peer) close() {
	// 에러가 있거나 채널이 닫혔을 때 peer을 닫을때 사용
	p.conn.Close()
	delete(Peers, p.key)
}

func (p *peer) read() {
	// delete peer in case of error
	defer p.close()
	for {
		_, m, err := p.conn.ReadMessage() // block
		if err != nil {
			break
		}
		fmt.Printf("%s", m)
	}
}

func (p *peer) write() {
	defer p.close()
	for {
		m, ok := <-p.inbox // block
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	key := fmt.Sprintf("%s:%s", address, port)

	p := &peer{
		conn:    conn,
		inbox:   make(chan []byte),
		key:     key,
		address: address,
		port:    port,
	}

	go p.read()
	go p.write()

	Peers[key] = p
	return p
}
