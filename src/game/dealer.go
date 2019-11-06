package game

import (
	"../cards"
	"fmt"
)

// Dealer is the blackjack dealer
type Dealer struct {
	Hand *cards.Hand
}

func NewDealer() *Dealer {
	return &Dealer{
		Hand: cards.NewHand(),
	}
}

func (dealer *Dealer) Move() string {
	value, _ := dealer.Hand.Value()
	if value < 17 {
		return "HIT"
	}
	return "STAY"
}

func (dealer *Dealer) Reset() {
	dealer.Hand = cards.NewHand()
}

func (dealer *Dealer) RevealCard() *cards.Card {
	return dealer.Hand.RevealCard()
}

func (dealer *Dealer) PrintHand(hasHuman bool) {
	if hasHuman {
		fmt.Printf("===== DEALER HAND =====\n")
		fmt.Printf("%s\n", dealer.Hand.Stringify())
	} else {
		fmt.Printf("Dealer has %s\n", dealer.Hand.Print())
	}
}
