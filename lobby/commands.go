package lobby

const (
	LOBBY_NEW  = "NEW"
	LOBBY_JOIN = "JOIN"
)

type LobbyCmd struct {
	Name   string `json:"name"`
	Action string `json:"action"`
	GameId string `json:"gameId"`
}

func (cmd *LobbyCmd) Valid() bool {
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
