package websocktoe

import (
	"log"
	"github.com/jakecoffman/websocktoe/random"
	"github.com/gorilla/websocket"
	"errors"
)

// dictates the front-end's view
const (
	VIEW_SETUP = "setup" // join or new?
	VIEW_PLAY = "game" // play
)

func loop(conn *websocket.Conn, games *Games) {
	player := &Player{id: random.PlayerId(), Name: "", conn: conn}
	defer func() {
		games.Disconnect(player)
	}()
	game := NewGame(player)

	for {
		err := conn.WriteJSON(game)
		if err != nil {
			log.Println(err)
			return
		}
		switch game.View {
		case VIEW_SETUP:
			// user needs to choose whether to join a game or create a new one
			if err = lobby(player, game, games); err != nil {
				log.Println(err)
				return
			}
		case VIEW_PLAY:
			if err = playGame(player, game); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

type lobbyCommand struct {
	Name   string `json:"name"`
	GameId string `json:"gameId"`
	Choice string `json:"choice"`
}

const (
	SETUP_NEW = "new"
	SETUP_JOIN = "join"
)

func lobby(player *Player, game *Game, games *Games) error {
	for {
		cmd := &lobbyCommand{}
		log.Println("Waiting for lobby command")
		err := player.conn.ReadJSON(cmd)
		log.Printf("Got lobby command: %#v\n", cmd)
		if err != nil {
			return err
		}
		if cmd.Name == "" {
			continue
		}
		player.Name = cmd.Name

		switch cmd.Choice {
		case SETUP_NEW:
			// user is trying a custom game ID
			if cmd.GameId != "" {
				// check to see if this game ID already exists
				if _, ok := games.Get(cmd.GameId); !ok {
					// successfully created new game with certain ID
					game.Id = cmd.GameId
					game.View = VIEW_PLAY
					games.Set(game.Id, game)
					log.Println("Starting a new game with custom ID")
					return nil
				} else {
					// game already exists with that ID
					// TODO: send error message?
					log.Println("User tried to start a new game but the game ID already exists")
					continue
				}
			}
			log.Println("Starting a new game")
			game.Id = random.GameId()
			game.View = VIEW_PLAY
			games.Set(game.Id, game)
			return nil
		case SETUP_JOIN:
			if cmd.GameId == "" {
				log.Println("User did not provide game ID")
				continue
			}
			if existingGame, ok := games.Get(cmd.GameId); !ok {
				log.Println("Existing game ID doesn't exist")
				continue
			} else {
				game = existingGame
				ok = game.Join(player)
				if !ok {
					log.Println("Could not join game")
					continue
				}
				return nil
			}
		default:
			continue
		}
	}
}

func playGame(player *Player, game *Game) error {
	return errors.New("Not yet implemented")
}
