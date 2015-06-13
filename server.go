package websocktoe

import (
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/websocket"
	"github.com/jakecoffman/websocktoe/tictactoe"
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
	games := tictactoe.NewPitBoss()
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
		log.Println(tictactoe.Loop(conn, games))
	})

	return mux
}
