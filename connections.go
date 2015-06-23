package websocktoe

import (
	"github.com/gorilla/websocket"
	"sync"
)

type none struct{}

type Connections struct {
	sync.RWMutex
	conns map[*websocket.Conn]none
}

func NewConnections() *Connections {
	return &Connections{conns: map[*websocket.Conn]none{}}
}

func (c *Connections) Add(conn *websocket.Conn) {
	c.Lock()
	defer c.Unlock()
	c.conns[conn] = none{}
}

func (c *Connections) Delete(conn *websocket.Conn) {
	c.Lock()
	defer c.Unlock()
	delete(c.conns, conn)
}

func (c *Connections) Size() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.conns)
}
