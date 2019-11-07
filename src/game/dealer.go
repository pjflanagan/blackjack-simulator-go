package game

import (
	"../cards"
	"fmt"
)

const (
	// HIT_ON_SOFT_17 is the rule for the dealer
	HIT_ON_SOFT_17 = true
)

// Dealer is the blackjack dealer
type Dealer struct {
	Hand *cards.Hand
}

// NewDealer returns a enw dealer with a hand
func NewDealer() *Dealer {
	return &Dealer{
		Hand: cards.NewHand(),
	}
}

// Deal adds a card to the dealer's hand
func (dealer *Dealer) Deal(card *cards.Card) {
	dealer.Hand.Add(card)
}

// Move returns a string representing the dealer's move
func (dealer *Dealer) Move() string {
	value, soft := dealer.Hand.Value()
	if value < 17 {
		return "HIT"
	} else if value == 17 && soft && HIT_ON_SOFT_17 {
		return "HIT"
	}
	return "STAY"
}

// DidBust returns true if the dealer has bust
func (dealer *Dealer) DidBust() bool {
	return dealer.Hand.DidBust()
}

// Reset gives the dealer a new hand
func (dealer *Dealer) Reset() {
	dealer.Hand = cards.NewHand()
}

// RevealCard returns the dealer's hidden card
func (dealer *Dealer) RevealCard() *cards.Card {
	return dealer.Hand.RevealCard()
}

// PrintHand prints the dealer's hand
func (dealer *Dealer) PrintHand(hasHuman bool) {
	if hasHuman {
		fmt.Printf("===== DEALER HAND =====\n")
		fmt.Printf("%s\n", dealer.Hand.LongformString())
	} else {
		fmt.Printf("Dealer has %s\n", dealer.Hand.ShorthandString())
	}
}
