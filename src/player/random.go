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
	randomPlayer := &RandomPlayer{
		basePlayer: initBasePlayer("Random"),
	}
	randomPlayer.basePlayer.child = randomPlayer
	return randomPlayer
}

// STEP 1: Bet -------------------------------------------------------------------------------------

// CanBet returns true when a player can bet
func (player *RandomPlayer) CanBet(minBet int) bool {
	return player.Chips > minBet && player.Status == c.PLAYER_READY
}

// Bet prompts a player to bet
func (player *RandomPlayer) Bet(minBet int, count int) {
	bet := minBet
	fmt.Printf("Random bets %d of %d chips available.\n", bet, player.Chips)
	player.bet(0, bet)
	return
}

// Blackjack handles when a player hits blackjack
func (player *RandomPlayer) Blackjack() {
	fmt.Printf("%s hit blackjack!\n", player.Name)
	player.blackjack()
}

// STEP 3: Move ------------------------------------------------------------------------------------

// Move returns string representing the move
func (player *RandomPlayer) Move(handIdx int) (move int) {
	fmt.Printf("%s has%s.\n", player.Name, player.Hands[handIdx].ShorthandString())
	validMoves := player.Hands[handIdx].GetValidMoves(player.Chips)
	move = validMoves[rand.Intn(len(validMoves))]
	return
}

// Hit returns true if hand is still active
func (player *RandomPlayer) Hit(handIdx int, card *cards.Card) bool {
	fmt.Printf("%s hits and receives %s.\n", player.Name, card.Stringify())
	return player.hit(handIdx, card)
}

// Split splits the player's hand
func (player *RandomPlayer) Split(handIdx int) {
	fmt.Printf("%s splits.\n", player.Name)
	player.split(handIdx)
}

// DoubleDown doubles down
func (player *RandomPlayer) DoubleDown(handIdx int, card *cards.Card) {
	fmt.Printf("%s doubles down and receives %s.\n", player.Name, card.Stringify())
	player.doubleDown(handIdx, card)
}

// Bust busts the players hand and sets the status
func (player *RandomPlayer) Bust(handIdx int) {
	fmt.Printf("%s busts and loses %d.\n", player.Name, player.Hands[handIdx].Wager)
	player.bust(handIdx)
}

// Stay returns true if the player's turn is still active
func (player *RandomPlayer) Stay(handIdx int) {
	fmt.Printf("%s stays.\n", player.Name)
	player.stay(handIdx)
}

// STEP 4: Payout ----------------------------------------------------------------------------------

// Payout print's message hand handles the payout
func (player *RandomPlayer) Payout(dealerHand *cards.Hand) {
	for i, hand := range player.Hands {
		result := hand.Result(dealerHand)

		switch result {
		case c.RESULT_BLACKJACK:
			// do not call payout for blackjack, money has already been given
			fmt.Printf("Random had a blackjack!\n")
		case c.RESULT_WIN:
			fmt.Printf("Random won!\n")
			player.payout(i, result)
		case c.RESULT_PUSH:
			fmt.Printf("Random pushes.\n")
			player.payout(i, result)
		case c.RESULT_BUST:
			// do not call payout for bust, money has already been taken
			fmt.Printf("Random busted.\n")
		case c.RESULT_LOSE:
			fmt.Printf("Random lost.\n")
			player.payout(i, result)
		}
	}

	if player.Chips > RANDOM_MAX_CHIPS {
		player.LeaveSeat()
	}
}

// HELPERS
