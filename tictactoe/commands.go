package tictactoe

type GameCmd struct {
	X     int  `json:"x"`
	Y     int  `json:"y"`
	Leave bool `json:"leave"`
}

func (cmd *GameCmd) Valid() bool {
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
