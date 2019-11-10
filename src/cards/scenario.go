package cards

type Scenario struct {
	HandString  string // string representing the player's hand
	UpcardValue int    // value of the dealer's upcard
}

func NewScenarioFromHands(hand *Hand, dealerHand *Hand, includePair bool) (Scenario, bool) {
	scenarioString := hand.StringScenarioCode(includePair)
	if scenarioString == "" {
		// TODO: start handling errors
		return Scenario{}, false
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
