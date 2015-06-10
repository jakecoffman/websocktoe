package websocktoe

import (
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/websocket"
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
	// TODO: Persistance
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
		Loop(player, games)
	})

	return mux
}

func Loop(player *Player, games *Games) {
	data := map[string]string{}
	err := player.ReadJSON(&data)
	if err != nil {
		log.Println(err)
		return
	}

	player.Name = data["name"]

	var game *Game
	switch data["action"] {
	case "NEW":
		game = games.NewGame(player)
	case "JOIN":
		game = games.Find(data["gameId"])
		if game == nil {
			log.Println("Game not found")
			return
		}
		ok := game.Join(player)
		if !ok {
			log.Println("error connecting")
			return
		}
	default:
		log.Println(data["action"])
		return
	}
	defer game.Leave(player)

	for {
		game.Update()

		// listen for messages from the player
		data := map[string]interface{}{}
		err := player.ReadJSON(&data)
		if err != nil {
			log.Println(err)
			return
		}

		// TODO: Handle commands here
	}
}
