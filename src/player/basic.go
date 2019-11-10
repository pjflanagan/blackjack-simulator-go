package player

import (
	"../cards"
	c "../constant"
	"../utils"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	BASIC_PAYOUT_LEAVE = 200
)

var SCENARIO_MOVE map[cards.Scenario]int

type BasicStrategyPlayer struct {
	basePlayer
}

// NewBasicStrategyPlayer returns a player that plays basic strategy
func NewBasicStrategyPlayer() *BasicStrategyPlayer {
	makeScenarioMoveMap()
	return &BasicStrategyPlayer{
		basePlayer: initBasePlayer("Basic"),
	}
}

func makeScenarioMoveMap() {
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
	default:
		return c.MOVE_STAY
	}
}

// Bet ----------------------------------------------------------------------------------------------

// CanBet returns true when a player can bet
func (player *BasicStrategyPlayer) CanBet(minBet int) bool {
	return player.Chips >= minBet && player.Status == c.PLAYER_READY
}

// Bet basic players bet the minumum
func (player *BasicStrategyPlayer) Bet(minBet int, count int) {
	bet := minBet
	fmt.Printf("%s bets the minimum %d of %d chips available.\n", player.Name, bet, player.Chips)
	player.bet(bet)
	return
}

// Move ------------------------------------------------------------------------------------

// Move returns string representing the move
func (player *BasicStrategyPlayer) Move(handIdx int, dealerHand *cards.Hand) (move int) {
	fmt.Printf("%s has %s.\n", player.Name, player.Hands[handIdx].StringSumReadable())
	validMoves := player.Hands[handIdx].GetValidMoves(player.Chips)
	if len(validMoves) == 0 {
		return c.MOVE_STAY
	}
	s, _ := cards.NewScenarioFromHands(player.Hands[0], dealerHand, true)
	move = SCENARIO_MOVE[s]
	if !utils.Contains(validMoves, move) {
		if move == c.MOVE_SPLIT && player.Chips < player.Hands[handIdx].Wager {
			s, _ := cards.NewScenarioFromHands(player.Hands[0], dealerHand, false)
			move = SCENARIO_MOVE[s]
		} else if move == c.MOVE_DOUBLE && player.Chips < player.Hands[handIdx].Wager {
			move = c.MOVE_HIT
		} else {
			fmt.Printf("\n ERROR, SHOULD NOT RETURN AN INVALID MOVE! \n")
		}
	}
	return
}

// Payout ----------------------------------------------------------------------------------

// Payout print's message hand handles the payout
func (player *BasicStrategyPlayer) Payout(dealerHand *cards.Hand) {
	for i, hand := range player.Hands {
		result := hand.Result(dealerHand)
		player.resultPayout(i, result)
	}

	if player.Chips > BASIC_PAYOUT_LEAVE {
		player.LeaveSeat()
	}
}

// HELPERS
