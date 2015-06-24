package tictactoe

type LobbyCommand struct {
	Name   string `json:"name"`
	Action string `json:"action"`
	GameId string `json:"gameId"`
}

func (cmd *LobbyCommand) Valid() bool {
	if cmd.Name == "" {
		return false
	}
	if cmd.Action == LOBBY_NEW {
		return true
	}
	if cmd.Action == LOBBY_JOIN && cmd.GameId != "" {
		return true
	}
	return false
}

type GameCommand struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Leave bool `json:"leave"`
}

func (cmd *GameCommand) Valid() bool {
	if cmd.Leave == true {
		return true
	}

	if cmd.X < 0 || cmd.X > 3 {
		return false
	}

	if cmd.Y < 0 || cmd.Y > 3 {
		return false
	}

	return true
}
