package websocktoe

import (
	"testing"
	"github.com/gorilla/websocket"
)

func TestConnections(t *testing.T) {
	conns := NewConnections()
	conn1 := &websocket.Conn{}
	conn2 := &websocket.Conn{}

	if conns.Size() != 0 {
		t.Error("Expected 0")
	}
	conns.Add(conn1)
	if conns.Size() != 1 {
		t.Error("Expected 1")
	}
	conns.Add(conn1)
	if conns.Size() != 1 {
		t.Error("Expected 1")
	}
	conns.Add(conn2)
	if conns.Size() != 2 {
		t.Error("Expected 2")
	}
	conns.Delete(conn1)
	if conns.Size() != 1 {
		t.Error("Expected 1")
	}
	conns.Delete(conn2)
	if conns.Size() != 0 {
		t.Error("Expected 0")
	}
}
