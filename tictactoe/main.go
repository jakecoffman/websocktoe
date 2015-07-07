package tictactoe

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/jakecoffman/websocktoe/lib"
)

const (
	LOBBY_NEW  = "NEW"
	LOBBY_JOIN = "JOIN"
)

func lobbyLoop(player *lib.Player, pitboss *PitBoss) (*Game, error) {
	var game *Game
	for {
		cmd := &LobbyCmd{}
		err := player.Read(cmd)
		if err != nil {
			return nil, err
		}
		if !cmd.Valid() {
			if err = player.Say("Invalid lobby command"); err != nil {
				return nil, err
			}
			log.Println(cmd)
			continue
		}

		player.Name = cmd.Name

		switch cmd.Action {
		case LOBBY_NEW:
			game = pitboss.NewGame(player)
		case LOBBY_JOIN:
			game = pitboss.Find(cmd.GameId)
			if game == nil {
				player.Say("Game not found")
				continue
			}
			ok := game.Join(player)
			if !ok {
				player.Say("error connecting")
				continue
			}
		default:
			player.Say(fmt.Sprintln("Unknown action, programmer error?", cmd.Action))
			continue
		}
		return game, nil
	}
}

func gameLoop(player *lib.Player, game *Game) error {
	for {
		cmd := &GameCmd{}
		err := player.Read(cmd)
		if err != nil {
			return err
		}
		if !cmd.Valid() {
			player.Say(fmt.Sprintf("Invalid command", cmd))
			continue
		}
		if cmd.Leave {
			return nil
		}
		if game.Over {
			player.Say("Game is over")
			continue
		}
		if !game.Move(player, cmd.X, cmd.Y) {
			player.Say("Invalid move")
			continue
		}

		game.Update()
	}
}

func Loop(conn *websocket.Conn, id string, pitboss *PitBoss) error {
	player, game := pitboss.RejoinOrNewPlayer(conn, id)
	defer player.Disconnect()
	var err error
	for {
		if game == nil {
			game, err = lobbyLoop(player, pitboss)
			if err != nil {
				log.Println(err)
				return err
			}
		} else {
			log.Printf("Player %v rejoining", player.Name)
		}
		game.Broadcast(fmt.Sprintf("Player %v has joined", player.Name))
		defer func() {
			game.Broadcast(fmt.Sprintf("Player %v has disconnected", player.Name))
			game.Update()
		}()

		game.Update()
		err = gameLoop(player, game)
		if err != nil {
			log.Println(err)
			return err
		}
		game.Leave(player)
		game.Broadcast(fmt.Sprintf("Player %v has left", player.Name))
		game.Update()
		game = nil
	}
}
