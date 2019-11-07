package cards

import (
	"fmt"
)

// Hand represents a single player's hand
type Hand struct {
	Cards        []*Card
	Wager        int
	HasBeenSplit bool
}

// NewHand initializes a new hand
func NewHand() *Hand {
	return new(Hand)
}

// Add appends a card to the hand
func (hand *Hand) Add(card *Card) {
	hand.Cards = append(hand.Cards, card)
}

// Split returns a second hand
func (hand *Hand) Split() (splitHand *Hand) {
	splitHand = new(Hand)
	splitHand.Cards = []*Card{hand.Cards[1]}
	splitHand.HasBeenSplit = true
	splitHand.Wager = hand.Wager
	hand.Cards = []*Card{hand.Cards[0]}
	hand.HasBeenSplit = true
	return
}

// RevealCard is used by the dealer to reveal all the cards
func (hand *Hand) RevealCard() *Card {
	hand.Cards[0].FlipUp()
	return hand.Cards[0]
}

// Value returns the value of the hand and accounts for aces
// if there is an ace it returns the higher value and if it is soft
func (hand *Hand) Value() (sum int, soft bool) {
	var aceCount int
	for _, card := range hand.Cards {
		// get the value of each card
		cardValue := card.Value()
		if cardValue == ACE_VALUE {
			aceCount += 1
		}
		// add the value to the sum
		sum += cardValue
	}
	sum = accountForAces(sum, aceCount)
	soft = aceCount > 0
	return
}

func accountForAces(sum int, aceCount int) int {
	if sum > 21 && aceCount > 0 {
		// if the value is over 21 but they have an ace then subtract 10
		// this ace is now being used as a 1 so subtract an ace from the aceCount
		return accountForAces(sum-10, aceCount-1)
	}
	// if they don't have any aces or the value is under 21 the existing sum is fine
	return sum
}

// DidBust returns true when the player has over 21
func (hand *Hand) DidBust() bool {
	value, _ := hand.Value()
	return didBust(value)
}

func didBust(value int) bool {
	return value > 21
}

// Is21 returns true when the player has 21
func (hand *Hand) Is21() bool {
	value, _ := hand.Value()
	return value == 21
}

// IsBlackjack returns true when the player has blackjack
func (hand *Hand) IsBlackjack() bool {
	value, _ := hand.Value()
	return hand.isBlackjack(value)
}

func (hand *Hand) isBlackjack(value int) bool {
	return value == 21 && len(hand.Cards) == 2 && !hand.HasBeenSplit
}

// Result looks at the player and dealer's hand and returns the string representing the result
// this should realy only be called when a player has STAYed and will only return WIN, PUSH, LOSE
func (hand *Hand) Result(dealerHand *Hand) string {
	playerHandValue, _ := hand.Value()
	dealerHandValue, _ := dealerHand.Value()
	switch {
	case didBust(playerHandValue):
		return "BUST"
	case hand.isBlackjack(playerHandValue):
		return "BLACKJACK"
	case didBust(dealerHandValue):
		return "WIN"
	case playerHandValue == dealerHandValue:
		return "PUSH"
	case playerHandValue > dealerHandValue:
		return "WIN"
	}
	return "LOSE"
}

// ShorthandString returns the shorthand
func (hand *Hand) ShorthandString() (str string) {
	for _, card := range hand.Cards {
		str = fmt.Sprintf("%s %s", str, card.Stringify())
	}
	return
}

// LongformString turns the hand into a string
func (hand *Hand) LongformString() (str string) {
	for row := 0; row < 7; row++ {
		for _, card := range hand.Cards {
			if row == 0 {
				str = fmt.Sprintf("%s┌─────────┐ ", str)
			} else if row == 6 {
				str = fmt.Sprintf("%s└─────────┘ ", str)
			} else if card.FaceDown {
				str = fmt.Sprintf("%s│░░░░░░░░░│ ", str)
			} else {
				cardinality, first, last := card.CardReadyStrings()
				switch row {
				case 1:
					str = fmt.Sprintf("%s│%s       │ ", str, first)
				case 2:
					str = fmt.Sprintf("%s│%s        │ ", str, cardinality)
				case 3:
					str = fmt.Sprintf("%s│         │ ", str)
				case 4:
					str = fmt.Sprintf("%s│        %s│ ", str, cardinality)
				case 5:
					str = fmt.Sprintf("%s│       %s│ ", str, last)
				}
			}
		}
		str = fmt.Sprintf("%s\n", str)
	}
	return
}
