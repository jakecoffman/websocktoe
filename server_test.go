package websocktoe

import (
	"testing"
	"github.com/gorilla/websocket"
	"net/http/httptest"
	"net/url"
	"net"
	"net/http"
	"log"
	"reflect"
)

func TestServer(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	server := httptest.NewServer(NewServer())
	defer server.Close()

	u, _ := url.Parse(server.URL)
	c1, err := net.Dial("tcp", u.Host)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	u.Path += "/ws"

	// http response ignored
	ws1, _, err := websocket.NewClient(c1, u, http.Header{"Origin": {server.URL}}, 1024, 1024)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer ws1.Close()

	type M map[string]interface{}

	err = ws1.WriteJSON(M{"name": "alice", "action": "NEW"})
	if err != nil {
		t.Fatal(err)
	}

	gameState := M{}
	err = ws1.ReadJSON(&gameState)
	if err != nil {
		t.Fatal(err)
	}

	c2, err := net.Dial("tcp", u.Host)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	ws2, _, err := websocket.NewClient(c2, u, http.Header{"Origin": {server.URL}}, 1024, 1024)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer ws2.Close()

	err = ws2.WriteJSON(M{"name": "bob", "action": "JOIN", "gameId": gameState["id"]})
	if err != nil {
		t.Fatal(err)
	}

	gameState = M{}
	err = ws1.ReadJSON(&gameState)
	if err != nil {
		t.Fatal(err)
	}

	gameState2 := M{}
	err = ws2.ReadJSON(&gameState2)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(gameState, gameState2) {
		t.Fatal("Game states not equal")
	}

//	out := M{"hello": "world"}
//	err = ws1.WriteJSON(out)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	msg := M{}
//	err = ws2.ReadJSON(&msg)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if !reflect.DeepEqual(msg, gameState) {
//		t.Fatalf("%#v != %#v", msg, gameState)
//	}
}
