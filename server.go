package websocktoe

import (
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"net/http/pprof"
)

var connections = 0

func New() *http.ServeMux {
	mux := http.NewServeMux()

	// TODO: remove for PROD
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	mux.Handle("/", http.FileServer(http.Dir("static")))

	upgrader := websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}

	// On server startup, create a new empty set of games
	// TODO: Persistance
	games := NewGames()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer func() {
			conn.Close()
			connections--
			log.Println("Connections", connections)
		}()
		connections++
		log.Println("Connections", connections)
		loop(conn, games)
	})

	return mux
}
