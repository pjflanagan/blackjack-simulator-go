package player

import (
	"../cards"
	c "../constant"
	"math/rand"
)

const (
// RANDOM_MAX_BET   = 30
)

// RandomPlayer extends basePlayer
type RandomPlayer struct {
	basePlayer
}

// NewRandomPlayer returns a new random player with name Random
func NewRandomPlayer(playerRules *PlayerRules) *RandomPlayer {
	return &RandomPlayer{
		basePlayer: initBasePlayer("Random", playerRules),
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

// HELPERS
