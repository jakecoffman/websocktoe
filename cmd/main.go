package main

import (
	"github.com/jakecoffman/websocktoe"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	mux := websocktoe.NewServer()
	hostport := "0.0.0.0:3030"
	log.Println("Running on", hostport)
	log.Fatal(http.ListenAndServe(hostport, mux))
}
