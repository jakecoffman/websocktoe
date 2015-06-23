package tictactoe

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

const (
	LOBBY_NEW  = "NEW"
	LOBBY_JOIN = "JOIN"
)

func Loop(conn *websocket.Conn, pitboss *PitBoss) error {
	player := NewPlayer(conn)
	var cmd *LobbyCommand
	for {
		cmd = &LobbyCommand{}
		err := player.ReadJSON(cmd)
		if err != nil {
			return err
		}
		if cmd.Valid() {
			break
		} else {
			log.Println("Invalid command:", cmd)
		}
	}

	player.Name = cmd.Name

	var game *Game
	switch cmd.Action {
	case LOBBY_NEW:
		game = pitboss.NewGame(player)
	case LOBBY_JOIN:
		game = pitboss.Find(cmd.GameId)
		if game == nil {
			return errors.New("Game not found")
		}
		ok := game.Join(player)
		if !ok {
			return errors.New("error connecting")
		}
	default:
		return errors.New(fmt.Sprintln("Unknown action, programmer error?", cmd.Action))
	}
	defer func() {
		game.Leave(player)
		log.Println(player.Name, "disconnected")
		game.Update()
	}()

	game.Update()
	for {
		cmd := &GameCommand{}
		err := player.ReadJSON(cmd)
		if err != nil {
			return err
		}
		if !cmd.Valid() {
			log.Println("Invalid command", cmd)
			continue
		}
		// TODO: add/check for leave command here
		if game.Over {
			log.Println("Game is over")
			continue
		}
		if !game.Move(player, cmd.X, cmd.Y) {
			log.Println("Invalid move")
			player.WriteJSON(struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			}{"message", "Invalid move"})
			continue
		}

		game.Update()
	}
}
