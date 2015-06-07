package websocktoe

import (
	"sync"
	"log"
)

type Game struct {
	sync.RWMutex
	Id      string `json:"id"`
	View    string `json:"view"`
	Players map[string]*Player `json:"players"`
}

// NewGame initializes a new game from scratch
func NewGame(player *Player) *Game {
	return &Game{Id: "", View: VIEW_SETUP, Players: map[string]*Player{player.id: player}}
}

// Join decides if a player may join a game or not, if so adds them as a player returns true, false otherwise
func (g *Game) Join(player *Player) bool {
	if len(g.Players) >= 2 {
		return false
	}
	g.Lock()
	defer g.Unlock()
	g.Players[player.id] = player

	// send update to all players that someone has joined!
	for _, player := range g.Players {
		err := player.conn.WriteJSON(g)
		if err != nil {
			log.Println("Player is no longer connected?", err)
		}
	}
	return true
}

func (g *Game) Leave(player *Player) {
	g.Lock()
	defer g.Unlock()
	delete(g.Players, player.id)
}

type Games struct {
	sync.RWMutex
	games map[string]*Game
}

func NewGames() *Games {
	return &Games{games: map[string]*Game{}}
}

func (g *Games) Get(id string) (*Game, bool) {
	g.RLock()
	defer g.RUnlock()
	game, ok := g.games[id]
	return game, ok
}

func (g *Games) Set(id string, game *Game) {
	g.Lock()
	defer g.Unlock()
	g.games[id] = game
}

func (g *Games) Disconnect(player *Player) {
	g.RLock()
	defer g.RUnlock()
	for _, game := range g.games {
		game.Leave(player)
	}
}
