package p2p

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Memory Safety를 위해 Mutex를 추가한 struct를 생성
type peers struct {
	v map[string]*peer
	m sync.Mutex // struct를 멀티 스레드의 접근으로부터 보호해줌
}

var Peers peers = peers{
	v: make(map[string]*peer),
}

type peer struct {
	conn    *websocket.Conn
	inbox   chan []byte
	key     string
	address string
	port    string
}

func AllPeers(p *peers) []string {
	p.m.Lock()
	defer p.m.Unlock()

	var keys []string
	for key := range p.v {
		keys = append(keys, key)
	}

	return keys
}

func (p *peer) close() {
	// 에러가 있거나 채널이 닫혔을 때 peer을 닫을때 사용
	Peers.m.Lock()
	defer Peers.m.Unlock()
	p.conn.Close()
	delete(Peers.v, p.key)
}

func (p *peer) read() {
	// delete peer in case of error
	defer p.close()
	for {
		m := Message{}
		err := p.conn.ReadJSON(&m) // block
		if err != nil {
			break
		}
		handleMsg(&m, p)
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
	Peers.m.Lock()
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

	Peers.v[key] = p
	Peers.m.Unlock()
	return p
}
