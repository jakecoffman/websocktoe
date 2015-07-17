package lib

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Game interface {
	Play(*Player, *PitBoss) error
	Id() string
	Leave(*Player)
	Find(string) *Player
	Broadcast(string, ...interface{})
	Update()
	Over() bool
	Join(*Player) bool
}

type PitBoss struct {
	sync.RWMutex
	games map[string]Game
}

func NewPitBoss() *PitBoss {
	return &PitBoss{games: map[string]Game{}}
}

// NewGame initializes a new game from scratch
func (g *PitBoss) Add(game Game) {
	g.Lock()
	defer g.Unlock()
	g.games[game.Id()] = game
}

func (g *PitBoss) Get(id string) (Game, bool) {
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

func (g *PitBoss) Find(gameId string) Game {
	g.RLock()
	defer g.RUnlock()
	game, ok := g.games[gameId]
	if !ok {
		return nil
	}
	return game
}

func (g *PitBoss) RejoinOrNewPlayer(conn *websocket.Conn, playerId string) (*Player, Game) {
	g.Lock()
	defer g.Unlock()
	for _, game := range g.games {
		if player := game.Find(playerId); player != nil {
			player.Rejoin(conn)
			return player, game
		}
	}
	return NewPlayer(conn, playerId), nil
}
