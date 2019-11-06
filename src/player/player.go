package player

import (
	"../cards"
)

// Player is the base class for all players (excluding dealer)
type Player interface {
	Move() string
	Payout(dealerHand *cards.Hand)
	// base
	HandString() string
	Deal(card *cards.Card)
	Reset()
}

type basePlayer struct {
	Hands      []*cards.Hand
	ActiveHand int
	Chips      int
}

func initBasePlayer() basePlayer {
	return basePlayer{
		Hands:      []*cards.Hand{cards.NewHand()},
		Chips:      100,
		ActiveHand: 0,
	}
}

// Reset resets the hand and the bet
func (player *basePlayer) Reset() {
	player.Hands = []*cards.Hand{cards.NewHand()}
	player.ActiveHand = 0
}

// Hand string calls the stringify function on the player's hand
func (player *basePlayer) HandString() string {
	return player.Hands[player.ActiveHand].Stringify()
}

// Deal adds a card to the player's hand
func (player *basePlayer) Deal(card *cards.Card) {
	player.Hands[player.ActiveHand].Add(card)
}

func (player *basePlayer) split() {
	splitHand := player.Hands[player.ActiveHand].Split()
	player.Hands = append(player.Hands, splitHand)
}

func (player *basePlayer) payout(handIdx int, result string) {
	wager := player.Hands[handIdx].Wager
	switch result {
	case "BLACKJACK":
		player.Chips += (wager * 3 / 2) + wager
	case "WIN":
		player.Chips += wager + wager
	case "PUSH":
		player.Chips += wager
	case "BUST", "LOSE":
	}
}
