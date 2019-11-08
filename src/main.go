package main

import (
	c "./constant"
	"./game"
)

func main() {
	// TODO: take command line input to determine what type of game we will be playing
	game := game.NewGame(20, 1)
	game.AddPlayer(c.TYPE_RANDOM)
	game.AddPlayer(c.TYPE_LEARNER)
	game.AddPlayer(c.TYPE_HUMAN)
	game.Play() // TODO: return a game.Summary object
}

// init game types
// GAME_COMPETITIVE: TYPE_BASIC, TYPE_RANDOM, TYPE_HUMAN
// GAME_LEARN: kick off a few go routines of learners idk
// GAME_HOUSE_ODDS: kick off a few go routines of TYPE_BASIC TYPE_COUNTER TYPE_RANDOM to compute house odds
