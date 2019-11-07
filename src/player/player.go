package player

import (
	"../cards"
	c "../constant"
)

// Player is the base class for all players (excluding dealer)
type Player interface {
	// Bet
	CanBet(minBet int) bool
	Bet(minBet int, count int)
	// Deal
	IsBlackjack() bool
	Blackjack()
	// Move
	Move(handIdx int) int
	Deal(handIdx int, card *cards.Card)
	Hit(handIdx int, card *cards.Card) bool
	Split(handIdx int)
	DoubleDown(handIdx int, card *cards.Card)
	Stay(handIDx int)
	// Payout
	Payout(dealerHand *cards.Hand)
	// Reset
	Reset(minBet int)
	LeaveSeat()
	// Helpers
	GetHands() []*cards.Hand
	StatusIs(statuses ...int) bool
}

type basePlayer struct {
	Name   string
	Hands  []*cards.Hand
	Chips  int
	Status int
}

func initBasePlayer(name string) basePlayer {
	return basePlayer{
		Name:   name,
		Hands:  []*cards.Hand{cards.NewHand()},
		Chips:  100,
		Status: c.PLAYER_READY,
	}
}

// STEP 1: Bet -------------------------------------------------------------------------------------

func (player *basePlayer) bet(handIdx int, bet int) {
	if player.Status == c.PLAYER_OUT {
		return
	}
	player.Chips -= bet
	player.Hands[handIdx].Wager = bet
	player.Status = c.PLAYER_ANTED
}

// STEP 2: Deal ------------------------------------------------------------------------------------

// Deal adds a card to the player's hand
func (player *basePlayer) Deal(handIdx int, card *cards.Card) {
	player.Hands[handIdx].Add(card)
	player.Status = c.PLAYER_JEPORADY
}

// IsBlackjack returns true when player gets a blackjack
func (player *basePlayer) IsBlackjack() bool {
	return player.Hands[0].IsBlackjack()
}

func (player *basePlayer) blackjack() {
	player.payout(0, c.RESULT_BLACKJACK)
	player.Status = c.PLAYER_BLACKJACK
}

// STEP 3: Turn ------------------------------------------------------------------------------------

func (player *basePlayer) validMoves() []string {
	return []string{}
}

// Hit returns true when hand is still active
func (player *basePlayer) hit(handIdx int, card *cards.Card) bool {
	player.Deal(handIdx, card)
	if player.Hands[handIdx].DidBust() {
		// if they bust then determine if turn is really over
		return player.bust(handIdx)
	} else if player.Hands[handIdx].Is21() {
		// if they hit 21 then this hand is over
		player.stay(handIdx)
		return false
	}
	// turn is still active, status is JEPORADY
	return true
}

func (player *basePlayer) split(handIdx int) {
	player.Chips -= player.Hands[handIdx].Wager
	splitHand := player.Hands[handIdx].Split()

	if handIdx == len(player.Hands)-1 {
		// if at the end of the array append the new hand to the end
		player.Hands = append(player.Hands, splitHand)
	} else {
		// if not at the end of the array then do something shifty
		// make space in the array for a new element
		player.Hands = append(player.Hands, nil)
		// copy over elements sourced from handIdx to one over
		copy(player.Hands[handIdx+2:], player.Hands[handIdx+1:])
		player.Hands[handIdx+1] = splitHand
	}

	player.Status = c.PLAYER_JEPORADY
}

func (player *basePlayer) doubleDown(handIdx int, card *cards.Card) {
	player.Hands[handIdx].Add(card)
	player.Chips -= player.Hands[handIdx].Wager
	player.Hands[handIdx].Wager *= 2
	if player.Hands[handIdx].DidBust() {
		player.bust(handIdx)
	} else {
		player.stay(handIdx)
	}
}

// Returns true if the player's hand is still active
func (player *basePlayer) bust(handIdx int) bool {
	if handIdx == len(player.Hands)-1 {
		player.payout(0, c.RESULT_BUST)
		player.Status = c.PLAYER_BUST
		return false
	}
	player.Status = c.PLAYER_JEPORADY
	return true
}

// Returns true if the player's turn is still active
func (player *basePlayer) stay(handIdx int) {
	if handIdx == len(player.Hands)-1 {
		player.Status = c.PLAYER_STAY
	} else {
		player.Status = c.PLAYER_JEPORADY
	}
}

// STEP 4: Payout ----------------------------------------------------------------------------------

func (player *basePlayer) payout(handIdx int, result int) {
	wager := player.Hands[handIdx].Wager
	player.Hands[handIdx].Wager = 0
	switch result {
	case c.RESULT_BLACKJACK:
		player.Chips += (wager * 3 / 2) + wager
	case c.RESULT_WIN:
		player.Chips += wager + wager
	case c.RESULT_PUSH:
		player.Chips += wager
	case c.RESULT_BUST, c.RESULT_LOSE:
	}
}

// STEP 5: Reset -----------------------------------------------------------------------------------

func (player *basePlayer) Reset(minBet int) {
	player.Hands = []*cards.Hand{cards.NewHand()}
	if player.Chips > minBet {
		player.Status = c.PLAYER_READY
	} else {
		player.Status = c.PLAYER_OUT
	}
}

// STEP 6: Leave -----------------------------------------------------------------------------------

func (player *HumanPlayer) LeaveSeat() {
	player.Status = c.PLAYER_OUT
}

// HELPERS -----------------------------------------------------------------------------------------

// Hands
func (player *basePlayer) GetHands() []*cards.Hand {
	return player.Hands
}

// StatusIs returns true if status is one of the strings
func (player *basePlayer) StatusIs(statuses ...int) bool {
	for _, status := range statuses {
		if status == player.Status {
			return true
		}
	}
	return false
}
