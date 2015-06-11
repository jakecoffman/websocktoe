package websocktoe

import (
	"sync"
	"log"
	"github.com/jakecoffman/websocktoe/random"
)

const (
	VIEW_SETUP = "SETUP"
	VIEW_PLAY = "PLAY"
)

type Game struct {
	sync.RWMutex
	Id       string `json:"id"`
	View     string `json:"view"`
	Players  map[string]*Player `json:"players"`
	Board    [3][3]string `json:"board"`
	LastMove string `json:"lastmove"` // name of the player that made the last move
	Winner   string `json:"winner"`
	Over     bool `json:"over"`
	Moves    int `json:"moves"`
}

func (g *Game) Join(player *Player) bool {
	g.Lock()
	defer g.Unlock()
	if len(g.Players) < 2 {
		g.Players[player.id] = player

		if len(g.Players) == 2 {
			g.View = VIEW_PLAY
		}
		return true
	}
	return false
}

func (g *Game) Leave(player *Player) {
	g.Lock()
	defer g.Unlock()
	delete(g.Players, player.id)
}

func (g *Game) Move(player *Player, x, y int) bool {
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
	for i := 0; i<3; i++ {
		if board[x][i] != name {
			break
		}
		if i == 2 {
			return name
		}
	}

	for i := 0; i<3; i++ {
		if board[i][y] != name {
			break
		}
		if i == 2 {
			return name
		}
	}

	for i := 0; i<3; i++ {
		if board[i][i] != name {
			break
		}
		if i == 2 {
			return name
		}
	}

	for i := 0; i<3; i++ {
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
	//	js, _ := json.MarshalIndent(g, "", "    ")
	//	log.Printf("UPDATE %v", string(js))
	for _, player := range g.Players {
		err := player.WriteJSON(g)
		if err != nil {
			log.Println(err)
		}
	}
}

type Games struct {
	sync.RWMutex
	games map[string]*Game
}

func NewGames() *Games {
	return &Games{games: map[string]*Game{}}
}

// NewGame initializes a new game from scratch
func (g *Games) NewGame(player *Player) *Game {
	g.Lock()
	defer g.Unlock()
	game := &Game{
		Id: random.GameId(),
		View: VIEW_PLAY,
		Players: map[string]*Player{player.id: player},
		Board: [3][3]string{},
		LastMove: "",
		Winner: "",
		Over: false,
	}
	game.Board[0] = [3]string{}
	game.Board[1] = [3]string{}
	game.Board[2] = [3]string{}
	g.games[game.Id] = game
	return game
}

func (g *Games) Get(id string) (*Game, bool) {
	g.RLock()
	defer g.RUnlock()
	game, ok := g.games[id]
	return game, ok
}

func (g *Games) Disconnect(player *Player) {
	g.RLock()
	defer g.RUnlock()
	for _, game := range g.games {
		game.Leave(player)
	}
}

func (g *Games) Find(gameId string) *Game {
	g.RLock()
	defer g.RUnlock()
	game, ok := g.games[gameId]
	if !ok {
		return nil
	}
	return game
}
