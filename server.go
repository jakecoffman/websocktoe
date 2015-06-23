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

	games := tictactoe.NewPitBoss()
	connections := NewConnections()

	mux.HandleFunc("/ws", HandlerTicTacToe(games, connections))
	return mux
}

func HandlerTicTacToe(games *tictactoe.PitBoss, conns *Connections) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		up := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, err := up.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer func() {
			conn.Close()
			conns.Delete(conn)
			log.Println("Connections", conns.Size())
		}()
		conns.Add(conn)
		log.Println("Connections", conns.Size())
		log.Println(tictactoe.Loop(conn, games))
	}
}
