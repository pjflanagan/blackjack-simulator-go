package game

import (
	"../cards"
	c "../constant"
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
func (dealer *Dealer) Move() int {
	value, handType := dealer.Hand.Value()
	if value < 17 {
		return c.MOVE_HIT
	} else if value == 17 && handType == c.HAND_SOFT && HIT_ON_SOFT_17 {
		return c.MOVE_HIT
	}
	return c.MOVE_STAY
}

// DidBust returns true if the dealer has bust
func (dealer *Dealer) DidBust() (didBust bool) {
	didBust = dealer.Hand.DidBust()
	if didBust {
		fmt.Printf("Dealer busts with%s.\n", dealer.Hand.ShorthandString())
	}
	return
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
		fmt.Printf("\n===== DEALER HAND =====\n")
		fmt.Printf("%s\n", dealer.Hand.LongformString())
	} else {
		fmt.Printf("Dealer has %s\n", dealer.Hand.ShorthandString())
	}
}
