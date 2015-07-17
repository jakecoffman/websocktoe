package tictactoe

import (
	"fmt"
	"log"

	"github.com/jakecoffman/websocktoe/lib"
)

func (g *TicTacToe) Play(player *lib.Player, pitboss *lib.PitBoss) error {
	//		if g == nil {
	//			tictactoe, err = lobbyLoop(player, pitboss)
	//			if err != nil {
	//				log.Println(err)
	//				return err
	//			}
	//		} else {
	//			log.Printf("Player %v rejoining", player.Name)
	//		}
	g.Broadcast("Player %v has joined", player.Name)
	g.Update()
	err := gameLoop(player, g)
	if err != nil {
		log.Println(err)
		return err
	}
	g.Leave(player)
	g.Broadcast(fmt.Sprintf("Player %v has left", player.Name))
	g.Update()
	return nil
}

func gameLoop(player *lib.Player, game *TicTacToe) error {
	for {
		cmd := &GameCmd{}
		err := player.Read(cmd)
		if err != nil {
			return err
		}
		if !cmd.Valid() {
			player.Say("Invalid command")
			continue
		}
		if cmd.Leave {
			return nil
		}
		if game.Over() {
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
