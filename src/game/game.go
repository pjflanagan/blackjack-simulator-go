package game

import ()

const ()

// Game represents a blackjack game
type Game struct {
	Table *Table
}

// NewGame returns a game with defaults
func NewGame(minBet int, deckCount int) *Game {
	return &Game{
		Table: NewTable(minBet, deckCount),
	}
}

// Play is the main game loop
func (game *Game) Play() {
	for i := 0; i < 4; i++ {
		game.Table.Deal()
		game.Table.Reset()
	}
}
