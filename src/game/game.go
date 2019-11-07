package game

import (
	"../player"
)

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

// AddPlayer adds a new player to the game
func (game *Game) AddPlayer(playerType string) {
	switch playerType {
	case "HUMAN":
		game.Table.TakeSeat(player.NewHumanPlayer(), true)
	}
}

// Play is the main game loop
func (game *Game) Play() {
	hasActivePlayer := true
	for hasActivePlayer {
		// while there are active players
		if hasActivePlayer = game.Table.TakeBets(); !hasActivePlayer {
			break
		}
		game.Table.Deal()
		game.Table.TakeTurns()
		game.Table.Payout()
		game.Table.Reset()
	}
	// game.Table.FinalStats()
}
