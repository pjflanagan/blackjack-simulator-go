package player

import (
	"../cards"
	c "../constant"
	"math"
)

type CardCounterPlayer struct {
	basePlayer
}

// NewCardCounterPlayer returns a player that plays basic strategy
func NewCardCounterPlayer(playerRules *PlayerRules) *CardCounterPlayer {
	makeScenarioMoveMap()
	return &CardCounterPlayer{
		basePlayer: initBasePlayer("Counter", playerRules),
	}
}

// CanBet returns true when a player can bet

// Bet based on true count
func (player *CardCounterPlayer) Bet(minBet int, trueCount float32) {
	var bet int
	if bettingUnits := trueCount - 1; bettingUnits > 1 {
		// if the betting unit is less greater than 1 then bet that
		bet = minBet * int(math.Floor(float64(trueCount)))
	} else {
		// otherwise bet the minimum
		bet = minBet
	}
	if bet > player.Chips {
		// if we chose to bet more than we have then bet it all
		bet = player.Chips
	}
	c.Print("%s bets %d of %d chips available.\n", player.Name, bet, player.Chips)
	player.bet(bet)
	return
}

// Move plays basic strategy
func (player *CardCounterPlayer) Move(handIdx int, dealerHand *cards.Hand) (move int) {
	c.Print("%s has %s.\n", player.Name, player.Hands[handIdx].StringSumReadable())
	return basicStrategyMove(player, handIdx, dealerHand)
}
