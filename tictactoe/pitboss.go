package tictactoe

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jakecoffman/websocktoe/lib"
	"github.com/jakecoffman/websocktoe/random"
)

type PitBoss struct {
	sync.RWMutex
	games map[string]*Game
}

func NewPitBoss() *PitBoss {
	return &PitBoss{games: map[string]*Game{}}
}

// NewGame initializes a new game from scratch
func (g *PitBoss) NewGame(player *lib.Player) *Game {
	g.Lock()
	defer g.Unlock()
	game := &Game{
		id:      random.GameId(),
		view:    VIEW_PLAY,
		players: map[string]*lib.Player{player.Id(): player},
		board:   [3][3]string{},
		over:    false,
	}
	game.board[0] = [3]string{}
	game.board[1] = [3]string{}
	game.board[2] = [3]string{}
	g.games[game.id] = game
	return game
}

func (g *PitBoss) Get(id string) (*Game, bool) {
	g.RLock()
	defer g.RUnlock()
	game, ok := g.games[id]
	return game, ok
}

func (g *PitBoss) Disconnect(player *lib.Player) {
	g.RLock()
	defer g.RUnlock()
	for _, game := range g.games {
		game.Leave(player)
	}
}

func (g *PitBoss) Find(gameId string) *Game {
	g.RLock()
	defer g.RUnlock()
	game, ok := g.games[gameId]
	if !ok {
		return nil
	}
	return game
}

func (g *PitBoss) RejoinOrNewPlayer(conn *websocket.Conn, playerId string) (*lib.Player, *Game) {
	g.Lock()
	defer g.Unlock()
	for _, game := range g.games {
		if player := game.Find(playerId); player != nil {
			player.Rejoin(conn)
			return player, game
		}
	}
	return lib.NewPlayer(conn, playerId), nil
}
