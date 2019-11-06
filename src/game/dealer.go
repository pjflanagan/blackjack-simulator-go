package game

import (
	"../cards"
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

func (dealer *Dealer) Reset() {
	dealer.Hand = cards.NewHand()
}
