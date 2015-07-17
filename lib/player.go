package lib

import (
	"errors"
	"sync"

	"fmt"
	"github.com/gorilla/websocket"
)

type Player struct {
	r         sync.Mutex
	w         sync.Mutex
	id        string
	Name      string `json:"name"`
	conn      *websocket.Conn
	Connected bool `json:"connected"`
}

func NewPlayer(conn *websocket.Conn, id string) *Player {
	return &Player{id: id, Name: "", conn: conn, Connected: true}
}

type PlayerMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (p *Player) Id() string {
	return p.id
}

func (p *Player) Rejoin(conn *websocket.Conn) {
	p.r.Lock()
	p.w.Lock()
	defer p.r.Unlock()
	defer p.w.Unlock()
	p.conn = conn
	p.Connected = true
}

func (p *Player) Say(message string, a ...interface{}) error {
	p.w.Lock()
	defer p.w.Unlock()
	if !p.Connected {
		return errors.New("Player disconnected")
	}
	err := p.conn.WriteJSON(PlayerMessage{"message", fmt.Sprintf(message, a...)})
	return err
}

func (p *Player) Read(data interface{}) error {
	p.r.Lock()
	defer p.r.Unlock()
	if !p.Connected {
		return errors.New("Player disconnected before read")
	}
	return p.conn.ReadJSON(data)
}

func (p *Player) Write(data interface{}) error {
	p.w.Lock()
	defer p.w.Unlock()
	if !p.Connected {
		return errors.New("Player disconnected")
	}
	return p.conn.WriteJSON(data)
}

func (p *Player) Disconnect() {
	p.w.Lock()
	p.r.Lock()
	defer p.w.Unlock()
	defer p.r.Unlock()
	p.Connected = false
}
