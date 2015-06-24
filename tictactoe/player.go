package tictactoe

import (
	"github.com/gorilla/websocket"
	"github.com/jakecoffman/websocktoe/random"
	"sync"
)

type Player struct {
	r    sync.Mutex
	w    sync.Mutex
	id   string
	Name string `json:"name"`
	conn *websocket.Conn
}

func NewPlayer(conn *websocket.Conn) *Player {
	return &Player{id: random.PlayerId(), Name: "", conn: conn}
}

type PlayerMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (p *Player) Say(message string) error {
	p.r.Lock()
	defer p.r.Unlock()
	err := p.conn.WriteJSON(PlayerMessage{"message", message})
	return err
}

func (p *Player) ReadLobbyCmd() (*LobbyCmd, error) {
	p.r.Lock()
	defer p.r.Unlock()
	lobbyCmd := &LobbyCmd{}
	err := p.conn.ReadJSON(lobbyCmd)
	return lobbyCmd, err
}

func (p *Player) ReadGameCmd() (*GameCmd, error) {
	p.r.Lock()
	defer p.r.Unlock()
	gameCmd := &GameCmd{}
	err := p.conn.ReadJSON(gameCmd)
	return gameCmd, err
}

func (p *Player) WriteGame(game *Game) error {
	p.w.Lock()
	defer p.w.Unlock()
	return p.conn.WriteJSON(struct {
		Type string `json:"type"`
		*Game
	}{Type: "state", Game: game})
}
