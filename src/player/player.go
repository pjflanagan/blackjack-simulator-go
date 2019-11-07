package player

import (
	"../cards"
	"fmt"
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
	Stay(handIDx int) bool
	// Payout
	Payout(dealerHand *cards.Hand)
	// Reset
	Reset(minBet int)
	LeaveSeat()
	// Helpers
	GetHands() []*cards.Hand
	StatusIs(statuses ...string) bool
	HandString(handIdx int) string
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

func (player *basePlayer) Blackjack() {
	fmt.Printf("%s hits blackjack!\n%s", player.Name, player.HandString(0))
	player.payout(0, "BLACKJACK")
	player.Status = "BLACKJACK"
}

// STEP 3: Turn ------------------------------------------------------------------------------------

// Hit returns true when turn is still active
func (player *basePlayer) Hit(handIdx int, card *cards.Card) bool {
	player.Deal(handIdx, card)
	fmt.Printf("%s receives %s.\n", player.Name, card.Stringify())
	if player.Hands[handIdx].DidBust() {
		return player.bust(handIdx)
	}
	return true
}

func (player *basePlayer) Split(handIdx int) {
	fmt.Printf("%s splits.\n", player.Name)
	player.Chips -= player.Hands[handIdx].Wager
	splitHand := player.Hands[handIdx].Split()
	player.Hands = append(player.Hands, splitHand)
	player.Status = "JEPORADY"
}

func (player *basePlayer) DoubleDown(handIdx int, card *cards.Card) {
	fmt.Printf("%s double down and receives %s.\n", player.Name, card.Stringify())
	player.Hands[handIdx].Add(card)
	player.Chips -= player.Hands[handIdx].Wager
	player.Hands[handIdx].Wager *= 2
	if player.Hands[handIdx].DidBust() {
		player.bust(handIdx)
	} else {
		player.Stay(handIdx)
	}
}

// Returns true if the player's turn is still active
func (player *basePlayer) bust(handIdx int) bool {
	fmt.Printf("%s busts and loses %d.\n", player.Name, player.Hands[handIdx].Wager)
	if handIdx == len(player.Hands)-1 {
		player.Status = "BUST"
		return false
	}
	player.Status = "JEPORADY"
	return true
}

// Returns true if the player's turn is still active
func (player *basePlayer) Stay(handIdx int) bool {
	fmt.Printf("%s stays.\n", player.Name)
	if handIdx == len(player.Hands)-1 {
		player.Status = "STAY"
		return false
	}
	player.Status = "JEPORADY"
	return true
}

// STEP 4: Payout ----------------------------------------------------------------------------------

func (player *basePlayer) Payout(dealerHand *cards.Hand) {
	for i, hand := range player.Hands {
		result := hand.Result(dealerHand)
		player.payout(i, result)
	}
}

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

// HandString calls the stringify function on the player's hand
func (player *basePlayer) HandString(handIdx int) string {
	return player.Hands[handIdx].LongformString()
}
