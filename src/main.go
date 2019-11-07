package main

import (
	c "./constant"
	"./game"
)

func main() {
	// TODO: take command line input to determine what type of game we will be playing
	game := game.NewGame(20, 1)
	game.AddPlayer(c.TYPE_HUMAN)
	game.Play()
}
