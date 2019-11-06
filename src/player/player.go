package player

import (
	"../cards"
)

// Player is the base class for all players (excluding dealer)
type Player interface {
	Move() string
}

type basePlayer struct {
	hand *cards.Hand
}
