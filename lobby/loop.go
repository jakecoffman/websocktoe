package lobby

import (
	"github.com/jakecoffman/websocktoe/lib"
	"github.com/jakecoffman/websocktoe/tictactoe"
)

func Loop(player *lib.Player, pitboss *lib.PitBoss) (lib.Game, error) {
	var game lib.Game
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
			continue
		}

		player.Name = cmd.Name

		switch cmd.Action {
		case LOBBY_NEW:
			game = tictactoe.NewTicTacToe(player)
			pitboss.Add(game)
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
			player.Say("Unknown action %#v", cmd.Action)
			continue
		}
		return game, nil
	}
}
