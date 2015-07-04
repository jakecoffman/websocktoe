package tictactoe

import (
	"github.com/gorilla/websocket"
	"sync"
	"errors"
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

func (p *Player) Say(message string) error {
	p.w.Lock()
	defer p.w.Unlock()
	if !p.Connected {
		return errors.New("Player disconnected")
	}
	err := p.conn.WriteJSON(PlayerMessage{"message", message})
	return err
}

func (p *Player) ReadLobbyCmd() (*LobbyCmd, error) {
	p.r.Lock()
	defer p.r.Unlock()
	if !p.Connected {
		return nil, errors.New("Player disconnected")
	}
	lobbyCmd := &LobbyCmd{}
	err := p.conn.ReadJSON(lobbyCmd)
	return lobbyCmd, err
}

func (p *Player) ReadGameCmd() (*GameCmd, error) {
	p.r.Lock()
	defer p.r.Unlock()
	if !p.Connected {
		return nil, errors.New("Player disconnected")
	}
	gameCmd := &GameCmd{}
	err := p.conn.ReadJSON(gameCmd)
	return gameCmd, err
}

func (p *Player) WriteGame(game *Game) error {
	p.w.Lock()
	defer p.w.Unlock()
	if !p.Connected {
		return errors.New("Player disconnected")
	}
	return p.conn.WriteJSON(struct {
		Type string `json:"type"`
		*Game
	}{Type: "state", Game: game})
}

func (p *Player) Disconnect() {
	p.w.Lock()
	p.r.Lock()
	defer p.w.Unlock()
	defer p.r.Unlock()
	p.Connected = false
}
