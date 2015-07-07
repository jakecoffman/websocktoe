package websocktoe

import (
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

type M map[string]interface{}

func TestServer(t *testing.T) {
	// setup test server
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	server := httptest.NewServer(NewServer())
	defer server.Close()

	url, _ := url.Parse(server.URL)
	url.Path += "/ws"

	// make first connection
	c1, err := net.Dial("tcp", url.Host)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}

	// prebake cookies and headers
	cookie1 := http.Cookie{Name: "ID", Value: "one"}
	cookie2 := http.Cookie{Name: "ID", Value: "two"}
	header1 := http.Header{"Origin": {server.URL}, "Cookie": {cookie1.String()}}
	header2 := http.Header{"Origin": {server.URL}, "Cookie": {cookie2.String()}}

	// create websocket connection (http response ignored)
	ws1, _, err := websocket.NewClient(c1, url, header1, 1024, 1024)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer ws1.Close()

	// write "new game" message to server
	err = ws1.WriteJSON(M{"name": "alice", "action": "NEW"})
	if err != nil {
		t.Fatal(err)
	}

	// read server response to new game
	gameState := M{}
	err = ws1.ReadJSON(&gameState)
	if err != nil {
		t.Fatal(err)
	}

	// create second connection to server and connect
	c2, err := net.Dial("tcp", url.Host)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	ws2, _, err := websocket.NewClient(c2, url, header2, 1024, 1024)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer ws2.Close()

	// join the game created by first connection
	err = ws2.WriteJSON(M{"name": "bob", "action": "JOIN", "gameId": gameState["id"]})
	if err != nil {
		t.Fatal(err)
	}

	// read both responses
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

	// check both responses are the same
	if !reflect.DeepEqual(gameState, gameState2) {
		t.Fatalf("Game states not equal: %#v %#v", gameState, gameState2)
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
