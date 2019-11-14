package player

import (
	"../cards"
	c "../constant"
	"../utils"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

var SCENARIO_MOVE map[cards.Scenario]int

type BasicStrategyPlayer struct {
	basePlayer
}

// NewBasicStrategyPlayer returns a player that plays basic strategy
func NewBasicStrategyPlayer(playerRules *PlayerRules) *BasicStrategyPlayer {
	makeScenarioMoveMap()
	return &BasicStrategyPlayer{
		basePlayer: initBasePlayer("Basic", playerRules),
	}
}

func makeScenarioMoveMap() {
	if len(SCENARIO_MOVE) != 0 {
		return
	}
	SCENARIO_MOVE = make(map[cards.Scenario]int)
	// Open the file
	csvfile, err := os.Open("./in/basic.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)

	// Iterate through the records
	for {
		// Read each record from csv
		handMoves, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for i := 1; i < len(handMoves); i++ {
			value := i
			if i == 1 {
				value = cards.ACE_VALUE
			}
			s := cards.NewScenario(handMoves[0], value)
			SCENARIO_MOVE[s] = getMoveFromString(handMoves[i])
		}
	}
}

func getMoveFromString(move string) int {
	switch strings.TrimSpace(move) {
	case "h":
		return c.MOVE_HIT
	case "p":
		return c.MOVE_SPLIT
	case "d":
		return c.MOVE_DOUBLE
	case "ds":
		return c.MOVE_DOUBLE_STAY
	default:
		return c.MOVE_STAY
	}
}

// Move ------------------------------------------------------------------------------------

// Move returns string representing the move
func (player *BasicStrategyPlayer) Move(handIdx int, dealerHand *cards.Hand) (move int) {
	c.Print("%s has %s.\n", player.Name, player.Hands[handIdx].StringSumReadable())
	return basicStrategyMove(player, handIdx, dealerHand)
}

func basicStrategyMove(player Player, handIdx int, dealerHand *cards.Hand) (move int) {
	hand := player.GetHand(handIdx)
	validMoves := hand.GetValidMoves(player.GetChips())
	if len(validMoves) == 0 {
		return c.MOVE_STAY
	}
	s, _ := cards.NewScenarioFromHands(player.GetHand(handIdx), dealerHand, true)
	move = SCENARIO_MOVE[s]
	if !utils.Contains(validMoves, move) {
		if move == c.MOVE_SPLIT {
			s, _ := cards.NewScenarioFromHands(player.GetHand(handIdx), dealerHand, false)
			move = SCENARIO_MOVE[s]
		} else if move == c.MOVE_DOUBLE {
			// if move is double but we don't have the funds then hit
			move = c.MOVE_HIT
		} else if move == c.MOVE_DOUBLE_STAY {
			// if double else stay but we don't have the funds then stay
			move = c.MOVE_STAY
		} else {
			log.Fatalf("\n Error, invalid move %d for scenario %+v", move, s)
		}
	}
	if move == c.MOVE_DOUBLE_STAY {
		// if the double was valid then actually double (nothing else recognized MOVE_DOUBLE_STAY)
		move = c.MOVE_DOUBLE
	}
	return
}

// HELPERS
