package websocktoe

import (
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/websocket"
	"github.com/jakecoffman/websocktoe/tictactoe"
	"github.com/nu7hatch/gouuid"
)

func NewServer() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	mux.Handle("/", AddSession(http.FileServer(http.Dir("static"))))

	games := tictactoe.NewPitBoss()
	connections := NewConnections()

	mux.HandleFunc("/ws", HandlerTicTacToe(games, connections))
	return mux
}

// AddSession ensures the connected user has a session
func AddSession(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("ID")
		if err != nil || cookie == nil || cookie.Value == "" {
			id, err := uuid.NewV4()
			if err != nil {
				log.Println("Error generating a UUID", err)
				w.WriteHeader(500)
				return
			}
			cookie = &http.Cookie{Name: "ID", Value: id.String()}
			http.SetCookie(w, cookie)
		}
		handler.ServeHTTP(w, r)
	}
}

func HandlerTicTacToe(games *tictactoe.PitBoss, conns *Connections) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("ID")
		if err != nil || cookie == nil || cookie.Value == "" {
			log.Println("No id in cookie", err)
			w.WriteHeader(500)
			return
		}
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
		log.Println(tictactoe.Loop(conn, cookie.Value, games))
	}
}
