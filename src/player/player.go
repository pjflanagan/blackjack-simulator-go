package player

import (
	"../cards"
)

// Player is the base class for all players (excluding dealer)
type Player interface {
	CanBet(minBet int) bool
	Bet(minBet int, count int) int
	Move(handIdx int) string
	Payout(dealerHand *cards.Hand)
	IsTurnOver(handIdx int) bool
	// base
	Deal(handIdx int, card *cards.Card)
	Reset()
	LeaveSeat()
	GetHands() []*cards.Hand
	IsActive() bool
	HandString(handIdx int) string
}

type basePlayer struct {
	Hands  []*cards.Hand
	Chips  int
	Active bool
}

func initBasePlayer() basePlayer {
	return basePlayer{
		Hands:  []*cards.Hand{cards.NewHand()},
		Chips:  100,
		Active: true,
	}
}

// STEP 1: Bet -------------------------------------------------------------------------------------

func (player *basePlayer) bet(bet int) int {
	player.Chips -= bet
	return bet
}

// STEP 2: Deal ------------------------------------------------------------------------------------

// Deal adds a card to the player's hand
func (player *basePlayer) Deal(handIdx int, card *cards.Card) {
	player.Hands[handIdx].Add(card)
}

// STEP 3: Turn ------------------------------------------------------------------------------------

func (player *basePlayer) split(handIdx int) {
	splitHand := player.Hands[handIdx].Split()
	player.Hands = append(player.Hands, splitHand)
}

// STEP 4: Payout ----------------------------------------------------------------------------------

func (player *basePlayer) Payout(dealerHand *cards.Hand) {
	for i, hand := range player.Hands {
		result := hand.Result(dealerHand)
		player.payout(i, result)
	}
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

// STEP 5: Reset -----------------------------------------------------------------------------------

func (player *basePlayer) Reset() {
	player.Hands = []*cards.Hand{cards.NewHand()}
}

// STEP 6: Leave -----------------------------------------------------------------------------------

func (player *HumanPlayer) LeaveSeat() {
	player.Active = false
}

// HELPERS -----------------------------------------------------------------------------------------

// Hands
func (player *basePlayer) GetHands() []*cards.Hand {
	return player.Hands
}

// IsActive
func (player *basePlayer) IsActive() bool {
	return player.Active
}

// HandString calls the stringify function on the player's hand
func (player *basePlayer) HandString(handIdx int) string {
	return player.Hands[handIdx].Stringify()
}

// IsTurnOver returns true when the turn is over
func (player *basePlayer) IsTurnOver(handIdx int) bool {
	return player.Hands[handIdx].DidBust() || player.Hands[handIdx].Is21()
}
