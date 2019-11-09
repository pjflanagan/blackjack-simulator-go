package cards

type Scenario struct {
	HandString  string // string representing the player's hand
	UpcardValue int    // value of the dealer's upcard
}

func NewScenario(hand *Hand, dealerHand *Hand) (Scenario, bool) {
	scenarioString := hand.StringScenarioCode()
	if scenarioString == "" {
		// TODO: start handling errors
		return Scenario{}, false
	}
	return Scenario{
		HandString:  scenarioString,
		UpcardValue: dealerHand.UpcardValue(),
	}, true
}