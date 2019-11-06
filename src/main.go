package main

import (
	"./game"
)

func main() {
	game := game.NewGame(20, 3)
	game.Play()
}
