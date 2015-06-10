package random

import (
	"math/rand"
	"time"
	"strconv"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

// TODO: Make these maintain a pool and only give unique

// GameId returns a random new Game ID (size 6)
func GameId() string {
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// PlayerId returns a new random Player ID
func PlayerId() string {
	// TODO: UUID?
	return strconv.Itoa(rand.Int())
}
