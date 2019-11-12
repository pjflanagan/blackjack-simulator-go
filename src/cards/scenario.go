package cards

import (
	c "../constant"
	"fmt"
	"log"
)

type Scenario struct {
	HandString  string // string representing the player's hand
	UpcardValue int    // value of the dealer's upcard
}

// NewScenarioFromHands returns true when the scenario is worth record
// that logic should maybe be moved
func NewScenarioFromHands(hand *Hand, dealerHand *Hand, includePair bool) (Scenario, bool) {
	scenarioString := hand.StringScenarioCode(includePair)
	switch scenarioString {
	case "bust", "blackjack":
		return Scenario{}, false
	case "":
		log.Fatal("NewScenarioFromHands: invalid scenario")
	}
	return Scenario{
		HandString:  scenarioString,
		UpcardValue: dealerHand.UpcardValue(),
	}, true
}

func NewScenario(handString string, upcardValue int) Scenario {
	return Scenario{
		HandString:  handString,
		UpcardValue: upcardValue,
	}
}

// StringScenarioCode returns the shorthand
// this belongs to the hand but is really only used in the context of scenarios
func (hand *Hand) StringScenarioCode(includePair bool) (str string) {
	value, handType := hand.Value()
	switch {
	case hand.isBlackjack(value):
		return "blackjack"
	case didBust(value):
		return "bust"
	case is21(value):
		return "21"
	case handType == c.HAND_PAIR:
		if hand.Cards[0].FaceName() == "A" {
			// this could also be called a soft12 but I don't wanna call it that because nobody ever would
			return "pairA"
		}
		if includePair {
			return fmt.Sprintf("pair%d", hand.Cards[0].Value())
		}
		return fmt.Sprintf("hard%d", value)
	case handType == c.HAND_SOFT:
		return fmt.Sprintf("soft%d", value)
	case handType == c.HAND_HARD:
		return fmt.Sprintf("hard%d", value)
	}
	return ""
}
