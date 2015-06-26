package tictactoe

import (
	"fmt"
	"github.com/gorilla/websocket"
)

const (
	LOBBY_NEW  = "NEW"
	LOBBY_JOIN = "JOIN"
)

func lobbyLoop(player *Player, pitboss *PitBoss) (*Game, error) {
	var game *Game
	for {
		lobbyCmd, err := player.ReadLobbyCmd()
		if err != nil {
			return nil, err
		}
		if !lobbyCmd.Valid() {
			if err = player.Say("Invalid command"); err != nil {
				return nil, err
			}
			continue
		}

		player.Name = lobbyCmd.Name

		switch lobbyCmd.Action {
		case LOBBY_NEW:
			game = pitboss.NewGame(player)
		case LOBBY_JOIN:
			game = pitboss.Find(lobbyCmd.GameId)
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
			player.Say(fmt.Sprintln("Unknown action, programmer error?", lobbyCmd.Action))
			continue
		}
		return game, nil
	}
}

func gameLoop(player *Player, game *Game) error {
	for {
		gameCmd, err := player.ReadGameCmd()
		if err != nil {
			return err
		}
		if !gameCmd.Valid() {
			player.Say(fmt.Sprintf("Invalid command", gameCmd))
			continue
		}
		if gameCmd.Leave {
			return nil
		}
		if game.Over {
			player.Say("Game is over")
			continue
		}
		if !game.Move(player, gameCmd.X, gameCmd.Y) {
			player.Say("Invalid move")
			continue
		}

		game.Update()
	}
}

func Loop(conn *websocket.Conn, pitboss *PitBoss) error {
	player := NewPlayer(conn)
	for {
		game, err := lobbyLoop(player, pitboss)
		if err != nil {
			return err
		}
		defer func() {
			game.Leave(player)
			game.Broadcast(fmt.Sprintf("Player %v has left", player.Name))
			game.Update()
		}()

		game.Update()
		err = gameLoop(player, game)
		if err != nil {
			return err
		}
		game.Leave(player)
		game.Broadcast(fmt.Sprintf("Player %v has left", player.Name))
		game.Update()
	}
}
