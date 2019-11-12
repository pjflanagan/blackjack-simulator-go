package player

import (
	"../cards"
	"fmt"
	"math"
)

const (
	COUNTER_PAYOUT_LEAVE = 300
)

type CardCounterPlayer struct {
	basePlayer
}

// NewCardCounterPlayer returns a player that plays basic strategy
func NewCardCounterPlayer() *CardCounterPlayer {
	makeScenarioMoveMap()
	return &CardCounterPlayer{
		basePlayer: initBasePlayer("Counter"),
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
	fmt.Printf("%s bets %d of %d chips available.\n", player.Name, bet, player.Chips)
	player.bet(bet)
	return
}

// Move plays basic strategy
func (player *CardCounterPlayer) Move(handIdx int, dealerHand *cards.Hand) (move int) {
	fmt.Printf("%s has %s.\n", player.Name, player.Hands[handIdx].StringSumReadable())
	return basicStrategyMove(player, handIdx, dealerHand)
}

// Payout print's message hand handles the payout
func (player *CardCounterPlayer) Payout(dealerHand *cards.Hand) {
	for i, hand := range player.Hands {
		result := hand.Result(dealerHand)
		player.resultPayout(i, result)
	}

	if player.Chips > COUNTER_PAYOUT_LEAVE {
		player.LeaveSeat()
	}
}
