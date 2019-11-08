package cards

import (
	"fmt"
)

const (
	ACE_VALUE  = 11
	FACE_VALUE = 10
)

// Card represents a playing card
type Card struct {
	// Face is int that represents the face of the card (not the value)
	Face        int
	cardinality int
	faceDown    bool
}

// NewCard returns a pointer to a new card
func NewCard(face int, cardinality int) *Card {
	return &Card{
		Face:        face,
		cardinality: cardinality,
		faceDown:    false,
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
	11: FACE_VALUE,
	12: FACE_VALUE,
	13: FACE_VALUE,
}

// Value returns the value of the card
// ace is 11 by default and requires additional
// checking in the Hand package
func (card *Card) Value() int {
	return GetFaceValue(card.Face)
}

func GetFaceValue(face int) int {
	return valueMap[face]
}

// FlipDown turns a card face down
func (card *Card) FlipDown() {
	card.faceDown = true
}

// IsFaceDown returns true when face is down
func (card *Card) IsFaceDown() bool {
	return card.faceDown
}

// FlipUp turns a card face up
func (card *Card) FlipUp() {
	card.faceDown = false
}

// IsFaceUp returns true when face is up
func (card *Card) IsFaceUp() bool {
	return !card.faceDown
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

// FaceName returns the name version of the face
func (card *Card) FaceName() string {
	return GetFaceString(card.Face)
}

func GetFaceString(face int) string {
	return faceMap[face]
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
	return fmt.Sprintf("%s%s", faceMap[card.Face], cardinalitySymbolMap[card.cardinality])
}

// CardReadyStrings returns cardinality, first face, last face
func (card *Card) CardReadyStrings() (string, string, string) {
	cardinality := cardinalitySymbolMap[card.cardinality]
	face := faceMap[card.Face]
	if len(face) == 2 {
		return cardinality, face, face
	}
	return cardinality, fmt.Sprintf("%s ", face), fmt.Sprintf(" %s", face)
}
