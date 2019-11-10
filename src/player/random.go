package player

import (
	"../cards"
	c "../constant"
	"fmt"
	"math/rand"
)

const (
	// RANDOM_MAX_BET   = 30
	RANDOM_MAX_CHIPS = 150
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

// CanBet returns true when a player can bet
func (player *RandomPlayer) CanBet(minBet int) bool {
	return player.Chips >= minBet && player.Status == c.PLAYER_READY
}

// Bet random players bet the minumum (TODO: make a random amount instead)
func (player *RandomPlayer) Bet(minBet int, count int) {
	bet := minBet
	fmt.Printf("%s bets %d of %d chips available.\n", player.Name, bet, player.Chips)
	player.bet(bet)
	return
}

// Move ------------------------------------------------------------------------------------

// Move returns string representing the move
func (player *RandomPlayer) Move(handIdx int, dealerHand *cards.Hand) (move int) {
	fmt.Printf("%s has %s.\n", player.Name, player.Hands[handIdx].StringSumReadable())
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
