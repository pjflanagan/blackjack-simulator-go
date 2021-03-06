package main

import (
	c "./constant"
	"./game"
	"./player"
	"./stats"
	"./utils"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	// "sync"
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

// STORY

var RANDOM_PLAYER_RULES_STORY = player.PlayerRules{
	LeavingChips: 150,
}

var BASIC_PLAYER_RULES_STORY = player.PlayerRules{
	LeavingChips: 200,
}

var COUNTER_PLAYER_RULES_STORY = player.PlayerRules{
	LeavingChips: 300,
}

func story(min int, decks int) {
	c.SetOutputMode(c.OUTPUT_LOG)
	game := game.NewGame(min, decks, 0)
	game.AddPlayer(c.TYPE_RANDOM, &RANDOM_PLAYER_RULES_STORY)
	game.AddPlayer(c.TYPE_BASIC, &BASIC_PLAYER_RULES_STORY)
	game.AddPlayer(c.TYPE_COUNTER, &COUNTER_PLAYER_RULES_STORY)
	game.Play()
}

// HUMAN

var HUMAN_PLAYER_RULES = player.PlayerRules{}

func human(min int, decks int) {
	c.SetOutputMode(c.OUTPUT_HUMAN)
	game := game.NewGame(min, decks, 0)
	game.AddPlayer(c.TYPE_RANDOM, &RANDOM_PLAYER_RULES_STORY)
	game.AddPlayer(c.TYPE_BASIC, &BASIC_PLAYER_RULES_STORY)
	game.AddPlayer(c.TYPE_COUNTER, &COUNTER_PLAYER_RULES_STORY)
	game.AddPlayer(c.TYPE_HUMAN, &HUMAN_PLAYER_RULES)
	game.Play()
}

// COMPARE

var RANDOM_PLAYER_RULES_COMPARE = player.PlayerRules{
	LeavingChips: 150,
	MaxHands:     15,
}

var BASIC_PLAYER_RULES_COMPARE = player.PlayerRules{
	LeavingChips: 200,
	MaxHands:     50,
}

var COUNTER_PLAYER_RULES_COMPARE = player.PlayerRules{
	LeavingChips: 300,
	MaxHands:     50,
}

func compare(min int, decks int) {
	var allStats []*stats.Stats
	for i := 0; i < 1000; i++ {
		c.SetOutputMode(c.OUTPUT_NONE)
		game := game.NewGame(min, decks, i)
		game.AddPlayer(c.TYPE_RANDOM, &RANDOM_PLAYER_RULES_COMPARE)
		game.AddPlayer(c.TYPE_BASIC, &BASIC_PLAYER_RULES_COMPARE)
		game.AddPlayer(c.TYPE_COUNTER, &COUNTER_PLAYER_RULES_COMPARE)
		gameStats := game.Play()
		allStats = append(allStats, gameStats...)
	}
	finalStats := stats.HouseOdds(allStats)
	for _, stat := range finalStats {
		fmt.Printf("%s strategy can expect a return of %f on each hand.\n", stat.GetStrategy(), stat.ExpectedGain())
	}
}

// LEARN

var LEARNER_PLAYER_RULES = player.PlayerRules{
	MaxHands: 10000,
}

func learn(min int, decks int) {
	c.SetOutputMode(c.OUTPUT_NONE)
	game := game.NewGame(min, decks, 0)
	game.AddPlayer(c.TYPE_LEARNER, &LEARNER_PLAYER_RULES)
	game.Play()
}
