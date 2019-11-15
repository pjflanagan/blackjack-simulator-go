package cards

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	// CARDINALITIES are 4
	CARDINALITIES = 4
	// FACES are 13
	FACES = 13
	// CARDS_IN_DECK
	CARDS_IN_DECK = 52
)

// Shoe has multiple decks of cards
type Shoe struct {
	deck            []*Card
	discards        []*Card
	reshuffleMarker int
	runNumber       int
}

// NewShoe returns a new shoe with a shuffled deck
func NewShoe(deckCount int, runNumber int) *Shoe {
	if deckCount < 1 {
		deckCount = 1
	}
	shoe := new(Shoe)
	shoe.runNumber = runNumber
	for decks := 0; decks < deckCount; decks++ {
		for face := 1; face <= FACES; face++ {
			for cardinality := 1; cardinality <= CARDINALITIES; cardinality++ {
				shoe.deck = append(shoe.deck, NewCard(face, cardinality))
			}
		}
	}
	shoe.Shuffle()
	shoe.setReshuffleMarker()
	return shoe
}

func (shoe *Shoe) setReshuffleMarker() {
	shoe.reshuffleMarker = len(shoe.deck) / 3
}

// Shuffle mkes the deck of cards random
func (shoe *Shoe) Shuffle() {
	// get all the cards back
	allCards := append(shoe.deck, shoe.discards...)
	// empty out the discards
	shoe.discards = []*Card{}
	// reset the deck
	shoe.deck = make([]*Card, len(allCards))
	// shuffle
	rand.Seed(time.Now().UnixNano() + int64(shoe.runNumber))
	randomDeckIdxs := rand.Perm(len(allCards))
	for i, deckIdx := range randomDeckIdxs {
		shoe.deck[deckIdx] = allCards[i]
		shoe.deck[deckIdx].FlipUp()
	}
	// shoe.cut()
	return
}

// func (shoe *Shoe) cut() {
// 	return
// }

// NeedsShuffle returns true when we are past the reshuffle marker
func (shoe *Shoe) NeedsShuffle() bool {
	return len(shoe.deck) < shoe.reshuffleMarker
}

// Take removes one from the top of the deck, holds onto it in the discard, and returns it
func (shoe *Shoe) Take() *Card {
	if len(shoe.deck) == 0 {
		fmt.Errorf("no cards left in deck")
	}
	// take from the top (back of deck)
	topCard := shoe.deck[len(shoe.deck)-1]
	// add this to the discards
	shoe.discards = append(shoe.discards, topCard)
	// remove this from the deck
	shoe.deck = shoe.deck[:len(shoe.deck)-1]
	// return this to the user
	return topCard
}

// Burn takes a card out of the shoe but does not return it (the table never sees it)
func (shoe *Shoe) Burn() {
	_ = shoe.Take()
}

func (shoe *Shoe) CardsRemaining() int {
	return len(shoe.deck)
}
