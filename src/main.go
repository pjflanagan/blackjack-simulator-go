package main

import (
	"./game"
)

func main() {
	// TODO: take command line input to determine what type of game we will be playing
	game := game.NewGame(20, 3)
	game.AddPlayer()
	game.Play()
}
