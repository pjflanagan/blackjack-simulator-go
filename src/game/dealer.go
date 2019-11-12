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

// Peek returns true when the dealer has blackjack
func (dealer *Dealer) Peek() bool {
	upcardValue := dealer.Hand.UpcardValue()
	if upcardValue == 10 || upcardValue == 11 {
		if dealer.Hand.IsBlackjack() {
			fmt.Printf("Dealer peeks and reveals blackjack.\n")
			return true
		}
		fmt.Printf("Dealer peeks and does not have blackjack.\n")
	}
	return false
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
func (dealer *Dealer) DidBust() bool {
	return dealer.Hand.DidBust()
}

// Reset gives the dealer a new hand
func (dealer *Dealer) Reset() {
	dealer.Hand = cards.NewHand()
}

// RevealCard returns the dealer's hidden card
func (dealer *Dealer) RevealCard() *cards.Card {
	revealCard := dealer.Hand.RevealCard()
	fmt.Printf("Dealer reveals %s. \n", revealCard.StringShorthand())
	return revealCard
}

// Hit adds card to the hand and outputs info
func (dealer *Dealer) Hit(card *cards.Card) {
	dealer.Deal(card)
	fmt.Printf("Dealer hits and receives %s.\n", card.StringShorthand())
}

// Stay prints that the dealer stays
func (dealer *Dealer) Stay() {
	fmt.Printf("Dealer stays.\n")
}

// PrintHand prints the dealer's hand
func (dealer *Dealer) PrintHand() {
	if c.OUTPUT_MODE == c.OUTPUT_NONE {
		return
	} else if c.OUTPUT_MODE == c.OUTPUT_HUMAN {
		fmt.Printf("\n===== Dealer Hand =====\n")
		fmt.Printf("%s\n", dealer.Hand.StringLongformReadable())
	} else {
		if dealer.Hand.Cards[1].IsFaceDown() {
			fmt.Printf("Dealer is showing %s.\n", dealer.Hand.StringShorthandReadable())
		} else if dealer.Hand.DidBust() {
			fmt.Printf("Dealer busts with %s.\n", dealer.Hand.StringShorthandReadable())
		} else {
			fmt.Printf("Dealer has %s.\n", dealer.Hand.StringShorthandReadable())
		}
	}
}
