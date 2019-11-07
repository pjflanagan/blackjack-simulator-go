package player

import (
	"../cards"
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
	Move(handIdx int) string
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
	StatusIs(statuses ...string) bool
}

type basePlayer struct {
	Name  string
	Hands []*cards.Hand
	Chips int
	// Status in order READY ANTED JEPORADY (BLACKJACK BUST STAY) OUT
	Status string
}

func initBasePlayer(name string) basePlayer {
	return basePlayer{
		Name:   name,
		Hands:  []*cards.Hand{cards.NewHand()},
		Chips:  100,
		Status: "READY",
	}
}

// STEP 1: Bet -------------------------------------------------------------------------------------

func (player *basePlayer) bet(handIdx int, bet int) {
	if player.Status == "LEFT" {
		return
	}
	player.Chips -= bet
	player.Hands[handIdx].Wager = bet
	player.Status = "ANTED"
}

// STEP 2: Deal ------------------------------------------------------------------------------------

// Deal adds a card to the player's hand
func (player *basePlayer) Deal(handIdx int, card *cards.Card) {
	player.Hands[handIdx].Add(card)
	player.Status = "JEPORADY"
}

// IsBlackjack returns true when player gets a blackjack
func (player *basePlayer) IsBlackjack() bool {
	return player.Hands[0].IsBlackjack()
}

func (player *basePlayer) blackjack() {
	player.payout(0, "BLACKJACK")
	player.Status = "BLACKJACK"
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
	// turn is still active, status is "JEPORADY"
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

	player.Status = "JEPORADY"
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
		player.payout(0, "BUST")
		player.Status = "BUST"
		return false
	}
	player.Status = "JEPORADY"
	return true
}

// Returns true if the player's turn is still active
func (player *basePlayer) stay(handIdx int) {
	if handIdx == len(player.Hands)-1 {
		player.Status = "STAY"
	} else {
		player.Status = "JEPORADY"
	}
}

// STEP 4: Payout ----------------------------------------------------------------------------------

func (player *basePlayer) payout(handIdx int, result string) {
	wager := player.Hands[handIdx].Wager
	player.Hands[handIdx].Wager = 0
	switch result {
	case "BLACKJACK":
		player.Chips += (wager * 3 / 2) + wager
	case "WIN":
		player.Chips += wager + wager
	case "PUSH":
		player.Chips += wager
	case "BUST", "LOSE":
	}
}

// STEP 5: Reset -----------------------------------------------------------------------------------

func (player *basePlayer) Reset(minBet int) {
	player.Hands = []*cards.Hand{cards.NewHand()}
	if player.Chips > minBet {
		player.Status = "READY"
	} else {
		player.Status = "OUT"
	}
}

// STEP 6: Leave -----------------------------------------------------------------------------------

func (player *HumanPlayer) LeaveSeat() {
	player.Status = "LEFT"
}

// HELPERS -----------------------------------------------------------------------------------------

// Hands
func (player *basePlayer) GetHands() []*cards.Hand {
	return player.Hands
}

// StatusIs returns true if status is one of the strings
func (player *basePlayer) StatusIs(statuses ...string) bool {
	for _, status := range statuses {
		if status == player.Status {
			return true
		}
	}
	return false
}
