package main

import (
	"math/rand"
	"time"

	"github.com/braheezy/hangman/internal"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	internal.Run()
}
