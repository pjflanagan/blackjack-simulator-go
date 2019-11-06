package cards

import (
	"fmt"
)

// Hand represents a single player's hand
type Hand struct {
	Cards []*Card
	Wager int
}

// NewHand initializes a new hand
func NewHand() *Hand {
	return new(Hand)
}

// Add appends a card to the hand
func (hand *Hand) Add(card *Card) {
	hand.Cards = append(hand.Cards, card)
}

// FlipUp is used by the dealer to reveal all the cards
func (hand *Hand) FlipUp() {
	for _, card := range hand.Cards {
		card.FlipUp()
	}
}

// Value returns the value of the hand and accounts for aces
// if there is an ace it returns the higher value and hasAce
func (hand *Hand) Value() (sum int, hasAce bool) {
	for _, card := range hand.Cards {
		// get the value of each card
		cardValue := card.Value()
		if cardValue == ACE_VALUE {
			// if it is an ace then store that to hasAce
			hasAce = true
		}
		// add the value to the sum
		sum += cardValue
	}
	if sum > 21 && hasAce {
		// if the value is over 21 but they have an ace then subtract 10
		sum -= 10
	}
	return
}

// Stringify turns the hand into a string
func (hand *Hand) Stringify() (str string) {
	for row := 0; row < 7; row++ {
		for _, card := range hand.Cards {
			if row == 0 {
				str = fmt.Sprintf("%s┌─────────┐", str)
			} else if row == 6 {
				str = fmt.Sprintf("%s└─────────┘", str)
			} else if card.FaceDown {
				str = fmt.Sprintf("%s│░░░░░░░░░│", str)
			} else {
				cardinality, first, last := card.CardReadyStrings()
				switch row {
				case 1:
					str = fmt.Sprintf("%s│%s       │", str, first)
				case 2:
					str = fmt.Sprintf("%s│%s        │", str, cardinality)
				case 3:
					str = fmt.Sprintf("%s│         │", str)
				case 4:
					str = fmt.Sprintf("%s│        %s│", str, cardinality)
				case 5:
					str = fmt.Sprintf("%s│       %s│", str, last)
				}
			}
		}
		str = fmt.Sprintf("%s\n", str)
	}
	return
}
