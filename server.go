package websocktoe

import (
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/websocket"
	"github.com/jakecoffman/websocktoe/lib"
	"github.com/jakecoffman/websocktoe/lobby"
	"github.com/nu7hatch/gouuid"
)

func NewServer() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	mux.Handle("/", AddSession(http.FileServer(http.Dir("static"))))

	games := lib.NewPitBoss()
	connections := NewConnections()

	mux.HandleFunc("/ws", WsHandler(games, connections))
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

func WsHandler(pitboss *lib.PitBoss, conns *Connections) http.HandlerFunc {
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
		player, game := pitboss.RejoinOrNewPlayer(conn, cookie.Value)
		defer func() {
			if game != nil {
				game.Broadcast("Player %v has disconnected", player.Name)
				game.Update()
			}
			player.Disconnect()
		}()
		if game == nil {
			game, err = lobby.Loop(player, pitboss)
			if err != nil {
				log.Println(err)
				return
			}
		}
		log.Println(game.Play(player, pitboss))
	}
}
