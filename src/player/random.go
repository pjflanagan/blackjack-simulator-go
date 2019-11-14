package player

import (
	"../cards"
	c "../constant"
	"math/rand"
)

const (
	// RANDOM_MAX_BET   = 30
	RANDOM_MAX_CHIPS = 150 // random only has to win 50 chips to stop playing
)

// RandomPlayer extends basePlayer
type RandomPlayer struct {
	basePlayer
}

// NewRandomPlayer returns a new random player with name Random
func NewRandomPlayer() *RandomPlayer {
	return &RandomPlayer{
		basePlayer: initBasePlayer("Random"),
	}
}

// Bet -------------------------------------------------------------------------------------

// Bet random players bet the minumum (TODO: make a random amount instead)
func (player *RandomPlayer) Bet(minBet int, trueCount float32) {
	bet := minBet
	c.Print("%s bets %d of %d chips available.\n", player.Name, bet, player.Chips)
	player.bet(bet)
	return
}

// Move ------------------------------------------------------------------------------------

// Move returns string representing the move
func (player *RandomPlayer) Move(handIdx int, dealerHand *cards.Hand) (move int) {
	c.Print("%s has %s.\n", player.Name, player.Hands[handIdx].StringSumReadable())
	validMoves := player.Hands[handIdx].GetValidMoves(player.Chips)
	if len(validMoves) == 0 {
		// this would happen if a player gets a 21 after a split (but we shouldn't go to here)
		move = c.MOVE_STAY
	} else {
		move = validMoves[rand.Intn(len(validMoves))]
	}
	return
}

// Payout ----------------------------------------------------------------------------------

// Payout print's message hand handles the payout
func (player *RandomPlayer) Payout(dealerHand *cards.Hand) {
	for i, hand := range player.Hands {
		result := hand.Result(dealerHand)
		player.resultPayout(i, result)
	}

	if player.Chips > RANDOM_MAX_CHIPS {
		player.LeaveSeat()
	}
}

// HELPERS
