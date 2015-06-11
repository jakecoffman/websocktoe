package websocktoe

import (
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/websocket"
	"errors"
	"fmt"
)

func NewServer() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	mux.Handle("/", http.FileServer(http.Dir("static")))

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// On server startup, create a new empty set of games
	games := NewGames()
	// all connections for broadcasting
	connections := map[*websocket.Conn]struct {}{}

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer func() {
			conn.Close()
			delete(connections, conn)
			log.Println("Connections", len(connections))
		}()
		connections[conn] = struct {}{}
		log.Println("Connections", len(connections))
		player := NewPlayer(conn)
		log.Println(Loop(player, games))
	})

	return mux
}

const (
	LOBBY_NEW = "NEW"
	LOBBY_JOIN = "JOIN"
)

func Loop(player *Player, games *Games) error {
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
		game = games.NewGame(player)
	case LOBBY_JOIN:
		game = games.Find(cmd.GameId)
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
	defer func(){
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
			continue
		}

		game.Update()
	}
}
