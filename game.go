package websocktoe

import (
	"sync"
	"log"
	"github.com/jakecoffman/websocktoe/random"
)

type Game struct {
	sync.RWMutex
	Id      string `json:"id"`
	View    string `json:"view"`
	Players map[string]*Player `json:"players"`
}

func (g *Game) Join(player *Player) bool {
	g.Lock()
	defer g.Unlock()
	if len(g.Players) < 2 {
		g.Players[player.id] = player
		return true
	}
	return false
}

func (g *Game) Leave(player *Player) {
	g.Lock()
	defer g.Unlock()
	delete(g.Players, player.id)
}

func (g *Game) Update() {
	g.RLock()
	defer g.RUnlock()
	for _, player := range g.Players {
		err := player.WriteJSON(g)
		if err != nil {
			log.Println(err)
		}
	}
}

func (g *Game) Say(player *Player, msg interface{}) {
	for id, player := range g.Players {
		if player.id != id {
			log.Println("Saying")
			player.WriteJSON(msg)
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
	game := &Game{
		Id: random.GameId(),
		View: "SETUP",
		Players: map[string]*Player{player.id: player},
	}
	g.Lock()
	defer g.Unlock()
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
