package main

import (
	"./cards"
	"fmt"
)

func main() {
	shoe := cards.NewShoe(1)
	var hands []*cards.Hand
	for i := 0; i < 4; i++ {
		hands = append(hands, cards.NewHand())
	}
	deal(shoe, hands)
}

// example, will be used by the table class
func deal(shoe *cards.Shoe, hands []*cards.Hand) {
	// burn a card
	_ = shoe.Take()
	for pass := 0; pass < 2; pass++ {
		for player, hand := range hands {
			card := shoe.Take()
			if player == len(hands)-1 && pass == 1 {
				// if this person is the dealer (should be last person)
				// and the card is the second pass
				card.FlipDown()
			}
			hand.Add(card)
			fmt.Printf("%s\n", hand.Stringify())
		}
	}
}
