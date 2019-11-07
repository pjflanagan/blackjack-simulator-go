package cards

import (
	"math/rand"
	"time"
)

const (
	// CARDINALITIES are 4
	CARDINALITIES = 4
	// FACES are 13
	FACES = 13
)

// Shoe has multiple decks of cards
type Shoe struct {
	Deck            []*Card
	discards        []*Card
	reshuffleMarker int
}

// NewShoe returns a new shoe with a shuffled deck
func NewShoe(deckCount int) *Shoe {
	shoe := new(Shoe)
	for decks := 0; decks < deckCount; decks++ {
		for face := 1; face <= FACES; face++ {
			for cardinality := 1; cardinality <= CARDINALITIES; cardinality++ {
				shoe.Deck = append(shoe.Deck, NewCard(face, cardinality))
			}
		}
	}
	shoe.Shuffle()
	return shoe
}

// Shuffle mkes the deck of cards random
func (shoe *Shoe) Shuffle() {
	// get all the cards back
	allCards := append(shoe.Deck, shoe.discards...)
	// empty out the discards
	shoe.discards = []*Card{}
	// reset the deck
	shoe.Deck = make([]*Card, len(allCards))
	// shuffle
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(allCards))
	for i, v := range perm {
		shoe.Deck[v] = allCards[i]
		shoe.Deck[v].FlipUp()
	}
	shoe.cut()
	shoe.setReshuffleMarker()
	return
}

func (shoe *Shoe) cut() {
	return
}

func (shoe *Shoe) setReshuffleMarker() {
	shoe.reshuffleMarker = len(shoe.Deck) / 3
}

// NeedsShuffle returns true when we are past the reshuffle marker
func (shoe *Shoe) NeedsShuffle() bool {
	return len(shoe.Deck) < shoe.reshuffleMarker
}

// Take removes one from the top of the deck, holds onto it in the discard, and returns it
func (shoe *Shoe) Take() *Card {
	// take from the top (back of deck)
	topCard := shoe.Deck[len(shoe.Deck)-1]
	// add this to the discards
	shoe.discards = append(shoe.discards, topCard)
	// remove this from the deck
	shoe.Deck = shoe.Deck[:len(shoe.Deck)-1]
	// return this to the user
	return topCard
}

func (shoe *Shoe) Burn() {
	_ = shoe.Take()
}
