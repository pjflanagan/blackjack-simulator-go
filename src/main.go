package main

import (
	c "./constant"
	"./game"
	"./utils"
	"log"
	"os"
	"reflect"
	"strconv"
)

var GAME_MODES = []string{"LEARN", "COMPARE", "COMPETE"}

func main() {
	mode, min, decks := readArgs()
	game := game.NewGame(min, decks)
	switch mode {
	case "LEARN":
		// TODO: a few go routines
		c.SetOutputMode(c.OUTPUT_NONE)
		game.AddPlayer(c.TYPE_LEARNER)
	case "COMPARE":
		// TODO: two versions of this, go routines that compute house odds for each,
		// and a version that just does one with a log for me to read
		game.AddPlayer(c.TYPE_RANDOM)
		game.AddPlayer(c.TYPE_BASIC)
		game.AddPlayer(c.TYPE_COUNTER)
	case "COMPETE":
		c.SetOutputMode(c.OUTPUT_HUMAN)
		game.AddPlayer(c.TYPE_RANDOM)
		game.AddPlayer(c.TYPE_BASIC)
		game.AddPlayer(c.TYPE_COUNTER)
		game.AddPlayer(c.TYPE_HUMAN)
	}
	// TODO: Play() should return a game.Summary object, so we can compile results here like odds and the scenario map
	game.Play()
}

func readArgs() (mode string, min int, decks int) {
	var err error
	// game mode
	if len(os.Args) > 1 && reflect.TypeOf(os.Args[1]).String() == "string" && utils.Contains(GAME_MODES, os.Args[1]) {
		mode = os.Args[1]
	} else {
		log.Fatalf("first arg must be one of %s\n", GAME_MODES)
	}
	// min bet
	if len(os.Args) > 2 && reflect.TypeOf(os.Args[2]).String() == "int" {
		if min, err = strconv.Atoi(os.Args[2]); err != nil || min <= 0 || min%2 != 0 {
			log.Fatalf("unable to parse valid min bet from second arg (must be positive even number) \n")
		}
	} else {
		min = c.DEFAULT_MIN
	}
	// deck count
	if len(os.Args) > 3 && reflect.TypeOf(os.Args[3]).String() == "int" {
		if decks, err = strconv.Atoi(os.Args[3]); err != nil || decks <= 0 {
			log.Fatalf("unable to parse valid deck count from third arg (must be positive number) \n")
		}
	} else {
		decks = c.DEFAULT_DECKS
	}

	return
}
