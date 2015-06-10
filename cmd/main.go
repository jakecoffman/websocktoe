package main

import (
	"log"
	"github.com/jakecoffman/websocktoe"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	mux := websocktoe.NewServer()
	log.Fatal(http.ListenAndServe("localhost:3030", mux))
}
