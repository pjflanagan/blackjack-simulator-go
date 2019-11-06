package game

import (
	"../cards"
	"../player"

	"fmt"
)

const ()

// Table represents a blackjack table
type Table struct {
	Shoe    *cards.Shoe
	Players []player.Player
	Dealer  *Dealer
	MinBet  int
}

// NewTable returns a table with defaults
func NewTable(minBet int, deckCount int) *Table {
	return &Table{
		Shoe:    cards.NewShoe(deckCount),
		Dealer:  NewDealer(),
		Players: []player.Player{player.NewHumanPlayer()},
		MinBet:  minBet,
	}
}

func (table *Table) Deal() {
	// burn a card
	table.Shoe.Burn()
	for pass := 0; pass < 2; pass++ {
		// make two passes for the deal
		for _, player := range table.Players {
			// deal each player in order
			card := table.Shoe.Take()
			player.Deal(card)
			fmt.Printf("%s\n", player.HandString())
		}
		// deal the dealer after the players
		card := table.Shoe.Take()
		if pass == 1 {
			// if it is the first pass keep if face down
			card.FlipDown()
		}
		table.Dealer.Hand.Add(card)
		fmt.Printf("%s\n", table.Dealer.Hand.Stringify())
	}
}

// Payout determines the winnings for each player
func (table *Table) Payout() {
	for _, player := range table.Players {
		player.Payout(table.Dealer.Hand)
	}
}

// remove
func (table *Table) Reset() {
	for _, player := range table.Players {
		player.Reset()
	}
	if table.Shoe.NeedsShuffle() {
		table.Shoe.Shuffle()
	}
	table.Dealer.Reset()
}
