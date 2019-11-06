package game

import (
	"../cards"
	"../player"
)

const ()

// Table represents a blackjack table
type Table struct {
	Shoe    *cards.Shoe
	Players []*player.Player
	Dealer  *Dealer
	MinBet  int
}

// NewTable returns a table with defaults
func NewTable(minBet int, deckCount int) *Table {
	return &Table{
		Shoe:   cards.NewShoe(deckCount),
		Dealer: new(Dealer),
		MinBet: minBet,
	}
}
