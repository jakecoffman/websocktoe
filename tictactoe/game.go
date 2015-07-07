package tictactoe

import (
	"log"
	"sync"

	"github.com/jakecoffman/websocktoe/lib"
)

const (
	VIEW_SETUP = "SETUP"
	VIEW_PLAY  = "PLAY"
)

type Game struct {
	sync.RWMutex
	Id       string                 `json:"id"`
	View     string                 `json:"view"`
	Players  map[string]*lib.Player `json:"players"`
	Board    [3][3]string           `json:"board"`
	LastMove string                 `json:"lastmove"` // name of the player that made the last move
	Winner   string                 `json:"winner"`
	Over     bool                   `json:"over"`
	Moves    int                    `json:"moves"`
}

func (g *Game) Join(player *lib.Player) bool {
	g.Lock()
	defer g.Unlock()
	if len(g.Players) < 2 {
		g.Players[player.Id()] = player

		if len(g.Players) == 2 {
			g.View = VIEW_PLAY
		}
		return true
	}
	return false
}

func (g *Game) Leave(player *lib.Player) {
	g.Lock()
	defer g.Unlock()
	delete(g.Players, player.Id())
}

func (g *Game) Broadcast(message string) {
	g.RLock()
	defer g.RUnlock()
	for _, p := range g.Players {
		_ = p.Say(message)
	}
}

func (g *Game) Move(player *lib.Player, x, y int) bool {
	g.Lock()
	defer g.Unlock()

	if g.LastMove == player.Name {
		return false
	}

	if g.Board[x][y] != "" {
		return false
	}

	g.Board[x][y] = player.Name
	g.LastMove = player.Name
	g.Moves++
	g.Winner = winner(g.Board, x, y, player.Name)
	if g.Winner != "" || g.Moves == 9 {
		g.Over = true
	}
	return true
}

func winner(board [3][3]string, x, y int, name string) string {
	for i := 0; i < 3; i++ {
		if board[x][i] != name {
			break
		}
		if i == 2 {
			return name
		}
	}

	for i := 0; i < 3; i++ {
		if board[i][y] != name {
			break
		}
		if i == 2 {
			return name
		}
	}

	for i := 0; i < 3; i++ {
		if board[i][i] != name {
			break
		}
		if i == 2 {
			return name
		}
	}

	for i := 0; i < 3; i++ {
		if board[i][2-i] != name {
			break
		}
		if i == 2 {
			return name
		}
	}

	return ""
}

func (g *Game) Update() {
	g.RLock()
	defer g.RUnlock()
	for _, player := range g.Players {
		game := struct {
			Type string `json:"type"`
			*Game
		}{Type: "state", Game: g}
		err := player.Write(game)
		if err != nil {
			log.Println(err)
		}
	}
}

func (g *Game) Find(playerId string) *lib.Player {
	g.RLock()
	defer g.RUnlock()
	player, _ := g.Players[playerId]
	return player
}
