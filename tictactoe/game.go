package tictactoe

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/jakecoffman/websocktoe/lib"
	"github.com/jakecoffman/websocktoe/random"
)

const (
	VIEW_SETUP = "SETUP"
	VIEW_PLAY  = "PLAY"
)

type TicTacToe struct {
	sync.RWMutex
	id       string
	view     string
	players  map[string]*lib.Player
	board    [3][3]string
	lastMove string // name of the player that went last
	winner   string
	over     bool
	moves    int
	messages []string
}

func NewTicTacToe(player *lib.Player) *TicTacToe {
	game := &TicTacToe{
		id:      random.GameId(),
		view:    VIEW_PLAY,
		players: map[string]*lib.Player{player.Id(): player},
		board:   [3][3]string{},
		over:    false,
	}
	game.board[0] = [3]string{}
	game.board[1] = [3]string{}
	game.board[2] = [3]string{}
	return game
}

func (g *TicTacToe) Id() string {
	g.RLock()
	defer g.RUnlock()
	return g.id
}

func (g *TicTacToe) Over() bool {
	g.RLock()
	defer g.RUnlock()
	return g.over
}

func (g *TicTacToe) Join(player *lib.Player) bool {
	g.Lock()
	defer g.Unlock()
	if len(g.players) < 2 {
		g.players[player.Id()] = player

		if len(g.players) == 2 {
			g.view = VIEW_PLAY
		}
		return true
	}
	return false
}

func (g *TicTacToe) Leave(player *lib.Player) {
	g.Lock()
	defer g.Unlock()
	delete(g.players, player.Id())
}

func (g *TicTacToe) Broadcast(message string, a ...interface{}) {
	g.Lock()
	defer g.Unlock()
	g.messages = append([]string{fmt.Sprintf(message, a...)}, g.messages...)
}

func (g *TicTacToe) Move(player *lib.Player, x, y int) bool {
	g.Lock()
	defer g.Unlock()

	if g.lastMove == player.Name {
		return false
	}

	if g.board[x][y] != "" {
		return false
	}

	g.board[x][y] = player.Name
	g.lastMove = player.Name
	g.moves++
	g.winner = winner(g.board, x, y, player.Name)
	if g.winner != "" || g.moves == 9 {
		g.over = true
	}
	return true
}

func (g *TicTacToe) Update() {
	g.RLock()
	defer g.RUnlock()
	for _, player := range g.players {
		err := player.Write(g)
		if err != nil {
			log.Println(err)
		}
	}
}

func (g *TicTacToe) Find(playerId string) *lib.Player {
	g.RLock()
	defer g.RUnlock()
	player, _ := g.players[playerId]
	return player
}

// MarshalJSON satisfies json.Marshaler interface
func (g *TicTacToe) MarshalJSON() ([]byte, error) {
	g.RLock()
	defer g.RUnlock()
	return json.Marshal(struct {
		Type     string                 `json:"type"`
		Id       string                 `json:"id"`
		View     string                 `json:"view"`
		Players  map[string]*lib.Player `json:"players"`
		Board    [3][3]string           `json:"board"`
		LastMove string                 `json:"lastmove"`
		Winner   string                 `json:"winner"`
		Over     bool                   `json:"over"`
		Moves    int                    `json:"moves"`
		Messages []string               `json:"messages"`
	}{
		"state",
		g.id,
		g.view,
		g.players,
		g.board,
		g.lastMove,
		g.winner,
		g.over,
		g.moves,
		g.messages,
	})
}

// UnmarshalJSON satisfies json.Unmarshaler interface
func (g *TicTacToe) UnmarshalJSON(data []byte) error {
	g.Lock()
	defer g.Unlock()
	return errors.New("We shant need this yet")
}
