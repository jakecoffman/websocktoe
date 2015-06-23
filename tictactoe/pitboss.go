package tictactoe

import (
	"github.com/jakecoffman/websocktoe/random"
	"sync"
)

type PitBoss struct {
	sync.RWMutex
	games map[string]*Game
}

func NewPitBoss() *PitBoss {
	return &PitBoss{games: map[string]*Game{}}
}

// NewGame initializes a new game from scratch
func (g *PitBoss) NewGame(player *Player) *Game {
	g.Lock()
	defer g.Unlock()
	game := &Game{
		Id:      random.GameId(),
		View:    VIEW_PLAY,
		Players: map[string]*Player{player.id: player},
		Board:   [3][3]string{},
		Over:    false,
	}
	game.Board[0] = [3]string{}
	game.Board[1] = [3]string{}
	game.Board[2] = [3]string{}
	g.games[game.Id] = game
	return game
}

func (g *PitBoss) Get(id string) (*Game, bool) {
	g.RLock()
	defer g.RUnlock()
	game, ok := g.games[id]
	return game, ok
}

func (g *PitBoss) Disconnect(player *Player) {
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
