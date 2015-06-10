package websocktoe

import (
	"github.com/gorilla/websocket"
	"github.com/jakecoffman/websocktoe/random"
	"sync"
)

type Player struct {
	r sync.RWMutex
	w sync.RWMutex
	id   string
	Name string `json:"name"`
	conn *websocket.Conn
}

func NewPlayer(conn *websocket.Conn) *Player {
	return &Player{id: random.PlayerId(), Name: "", conn: conn}
}

func (p *Player) ReadJSON(v interface{}) error {
	p.r.Lock()
	defer p.r.Unlock()
	err := p.conn.ReadJSON(&v)
	return err
}

func (p *Player) WriteJSON(v interface{}) error {
	p.w.Lock()
	defer p.w.Unlock()
	return p.conn.WriteJSON(v)
}
