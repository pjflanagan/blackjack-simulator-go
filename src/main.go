package main

import (
	c "./constant"
	"./game"
	"./stats"
	"./utils"
	"log"
	"os"
	"reflect"
	"strconv"
)

var GAME_MODES = []string{"LEARN", "COMPARE", "HUMAN", "STORY"}

func main() {
	mode, min, decks := readArgs()
	switch mode {
	case "LEARN":
		learn(min, decks)
	case "STORY":
		story(min, decks)
	case "COMPARE":
		compare(min, decks)
	case "HUMAN":
		human(min, decks)
	}
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

func story(min int, decks int) {
	c.SetOutputMode(c.OUTPUT_LOG)
	game := game.NewGame(min, decks)
	game.AddPlayer(c.TYPE_RANDOM)
	game.AddPlayer(c.TYPE_BASIC)
	game.AddPlayer(c.TYPE_COUNTER)
	game.Play()
}

func human(min int, decks int) {
	c.SetOutputMode(c.OUTPUT_HUMAN)
	game := game.NewGame(min, decks)
	game.AddPlayer(c.TYPE_RANDOM)
	game.AddPlayer(c.TYPE_BASIC)
	game.AddPlayer(c.TYPE_COUNTER)
	game.AddPlayer(c.TYPE_HUMAN)
	game.Play()
}

func compare(min int, decks int) {
	for i := 0; i < 20; i++ {
		go func() []*stats.Stats {
			c.SetOutputMode(c.OUTPUT_NONE)
			game := game.NewGame(min, decks)
			game.AddPlayer(c.TYPE_LEARNER)
			return game.Play()
		}()
	}
}

func learn(min int, decks int) {
	c.SetOutputMode(c.OUTPUT_NONE)
	game := game.NewGame(min, decks)
	game.AddPlayer(c.TYPE_LEARNER)
	game.Play()
}
