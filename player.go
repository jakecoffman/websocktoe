package websocktoe

import "github.com/gorilla/websocket"

type Player struct {
	id   string
	Name string `json:"name"`
	conn *websocket.Conn
}
