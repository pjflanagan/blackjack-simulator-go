package player

import (
	"../cards"
)

type HumanPlayer struct {
	basePlayer
}

func NewHumanPlayer() *HumanPlayer {
	return &HumanPlayer{
		basePlayer: initBasePlayer(),
	}
}

func (player *HumanPlayer) Move() string {
	return "HIT"
}

func (player *HumanPlayer) Payout(dealerHand *cards.Hand) {
	for i, hand := range player.Hands {
		result := hand.Result(dealerHand)
		player.payout(i, result)
	}
}
