package cards

import (
	c "../constant"
	"fmt"
)

// Hand represents a single player's hand
type Hand struct {
	Cards        []*Card
	Wager        int
	hasBeenSplit bool
}

// NewHand initializes a new hand
func NewHand() *Hand {
	return new(Hand)
}

// DEAL --------------------------------------------------------------------------------------------

// Add appends a card to the hand
func (hand *Hand) Add(card *Card) {
	hand.Cards = append(hand.Cards, card)
}

// Split returns a second hand
func (hand *Hand) Split() (splitHand *Hand) {
	splitHand = new(Hand)
	splitHand.Cards = []*Card{hand.Cards[1]}
	splitHand.hasBeenSplit = true
	splitHand.Wager = hand.Wager
	hand.Cards = []*Card{hand.Cards[0]}
	hand.hasBeenSplit = true
	return
}

// RevealCard is used by the dealer to reveal all the cards
// the second card is the hidden one
func (hand *Hand) RevealCard() *Card {
	hand.Cards[1].FlipUp()
	return hand.Cards[1]
}

// UpcardFace the face of the showing hand
// the first card is the showing one
func (hand *Hand) UpcardFace() int {
	return hand.Cards[0].Face
}

// UpcardValue the face of the showing hand
// the first card is the showing one
func (hand *Hand) UpcardValue() int {
	return hand.Cards[0].Value()
}

// VALUE -------------------------------------------------------------------------------------------

// Value returns the value of the hand and accounts for aces
// if there is an ace it returns the higher value and if it is soft
func (hand *Hand) Value() (int, int) {
	var aceCount int
	var sum int
	for _, card := range hand.Cards {
		// get the value of each card
		cardValue := card.Value()
		if cardValue == ACE_VALUE {
			aceCount++
		}
		// add the value to the sum
		sum += cardValue
	}

	value, isSoft := accountForAces(sum, aceCount)
	if hand.isPair() {
		return value, c.HAND_PAIR
	} else if isSoft {
		return value, c.HAND_SOFT
	} else {
		return value, c.HAND_HARD
	}
}

func accountForAces(sum int, aceCount int) (int, bool) {
	if sum > 21 && aceCount > 0 {
		// if the value is over 21 but they have an ace then subtract 10
		// this ace is now being used as a 1 so subtract an ace from the aceCount
		return accountForAces(sum-10, aceCount-1)
	}
	// if they don't have any aces or the value is under 21 the existing sum is fine
	return sum, aceCount > 0
}

func (hand *Hand) isPair() bool {
	// if there are two cards in the hand and the value's are the same
	return len(hand.Cards) == 2 && hand.Cards[0].Value() == hand.Cards[1].Value()
}

// MOVES -------------------------------------------------------------------------------------------

func (hand *Hand) canDoubleDown(chips int) bool {
	return len(hand.Cards) == 2 && chips > hand.Wager
}

// GetValidMoves returns an array of valid move ints
func (hand *Hand) GetValidMoves(chips int) (moves []int) {
	value, _ := hand.Value()
	if didBust(value) || is21(value) {
		return
	}
	moves = append(moves, c.MOVE_STAY, c.MOVE_HIT)
	if hand.isPair() && chips > hand.Wager {
		moves = append(moves, c.MOVE_SPLIT)
	}
	if hand.canDoubleDown(chips) {
		moves = append(moves, c.MOVE_DOUBLE)
	}
	return
}

// DidBust returns true when the player has over 21
func (hand *Hand) DidBust() bool {
	value, _ := hand.Value()
	return didBust(value)
}

func didBust(value int) bool {
	return value > 21
}

// RESULTS -----------------------------------------------------------------------------------------

// Result looks at the player and dealer's hand and returns the string representing the result
// this should realy only be called when a player has STAYed and will only return WIN, PUSH, LOSE
func (hand *Hand) Result(dealerHand *Hand) int {
	playerHandValue, _ := hand.Value()
	dealerHandValue, _ := dealerHand.Value()
	playerIsBlackjack := hand.isBlackjack(playerHandValue)
	switch {
	case dealerHand.isBlackjack(dealerHandValue) && playerIsBlackjack:
		return c.RESULT_PUSH
	case didBust(playerHandValue):
		return c.RESULT_BUST
	case playerIsBlackjack:
		return c.RESULT_BLACKJACK
	case didBust(dealerHandValue):
		return c.RESULT_WIN
	case playerHandValue == dealerHandValue:
		return c.RESULT_PUSH
	case playerHandValue > dealerHandValue:
		return c.RESULT_WIN
	}
	return c.RESULT_LOSE
}

// Is21 returns true when the player has 21
func (hand *Hand) Is21() bool {
	value, _ := hand.Value()
	return is21(value)
}

func is21(value int) bool {
	return value == 21
}

// IsBlackjack returns true when the player has blackjack
func (hand *Hand) IsBlackjack() bool {
	value, _ := hand.Value()
	return hand.isBlackjack(value)
}

func (hand *Hand) isBlackjack(value int) bool {
	return value == 21 && len(hand.Cards) == 2 && !hand.hasBeenSplit
}

// OUTPUT ------------------------------------------------------------------------------------------

// StringSumReadable returns the shorthand
func (hand *Hand) StringSumReadable() (str string) {
	value, handType := hand.Value()
	if didBust(value) {
		return "bust"
	} else if hand.isBlackjack(value) {
		return "blackjack"
	} else if handType == c.HAND_PAIR {
		if hand.Cards[0].FaceName() == "A" {
			return fmt.Sprintf("pair of aces")
		}
		return fmt.Sprintf("pair of %d's", value/2)
	} else if is21(value) {
		return "21"
	} else if handType == c.HAND_SOFT {
		return fmt.Sprintf("soft %d", value)
	}
	return fmt.Sprintf("hard %d", value)
}

// StringShorthandReadable returns the shorthand
func (hand *Hand) StringShorthandReadable() (str string) {
	if hand.Cards[1].IsFaceDown() {
		// if the second card is face down then only show the first card
		return hand.Cards[0].StringShorthand()
	}
	for i, card := range hand.Cards {
		if i == 0 {
			str = card.StringShorthand()
		} else {
			str = fmt.Sprintf("%s %s", str, card.StringShorthand())
		}
	}
	return
}

// StringLongformReadable turns the hand into a string
func (hand *Hand) StringLongformReadable() (str string) {
	for row := 0; row < 7; row++ {
		for _, card := range hand.Cards {
			if row == 0 {
				str = fmt.Sprintf("%s┌─────────┐ ", str)
			} else if row == 6 {
				str = fmt.Sprintf("%s└─────────┘ ", str)
			} else if card.IsFaceDown() {
				str = fmt.Sprintf("%s│░░░░░░░░░│ ", str)
			} else {
				cardinality, first, last := card.StringsForLongform()
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
