package cards

import (
	"fmt"
)

const (
	ACE_VALUE = 11
)

// Card represents a playing card
type Card struct {
	// Face is integer that represents the face value in FaceMap
	Face int
	// Cardinality is integer that represents the cardinality in CardinalityMap
	Cardinality int
	// FaceDown is bool representing if the card is facing down
	FaceDown bool
	// Do not store value in the card, the card's value depends on the situation
}

// NewCard returns a pointer to a new card
func NewCard(face int, cardinality int) *Card {
	return &Card{
		Face:        face,
		Cardinality: cardinality,
		FaceDown:    false,
	}
}

// valueMap maps integers to a value
var valueMap = map[int]int{
	1:  ACE_VALUE,
	2:  2,
	3:  3,
	4:  4,
	5:  5,
	6:  6,
	7:  7,
	8:  8,
	9:  9,
	10: 10,
	11: 10,
	12: 10,
	13: 10,
}

// Value returns the value of the card
// ace is 11 by default and requires additional
// checking in the Hand package
func (card *Card) Value() int {
	return valueMap[card.Face]
}

// FlipDown turns a card face down
func (card *Card) FlipDown() {
	card.FaceDown = true
}

// FlipUp turns a card face up
func (card *Card) FlipUp() {
	card.FaceDown = false
}

// faceMap maps integers to face strings
var faceMap = map[int]string{
	1:  "A",
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "10",
	11: "J",
	12: "Q",
	13: "K",
}

// cardinalitySymbolMap maps integers to a character
var cardinalitySymbolMap = map[int]string{
	1: "♠",
	2: "♣",
	3: "♥",
	4: "♦",
}

// Stringify returns a string with just the value and cardinality symbol
func (card *Card) Stringify() string {
	return fmt.Sprintf("%s%s", faceMap[card.Face], cardinalitySymbolMap[card.Cardinality])
}

// CardReadyStrings returns cardinality, first face, last face
func (card *Card) CardReadyStrings() (string, string, string) {
	cardinality := cardinalitySymbolMap[card.Cardinality]
	face := faceMap[card.Face]
	if len(face) == 2 {
		return cardinality, face, face
	}
	return cardinality, fmt.Sprintf("%s ", face), fmt.Sprintf(" %s", face)
}
