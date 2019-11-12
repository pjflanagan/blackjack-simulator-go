package game

import (
	c "../constant"
	"../player"
)

const ()

// type Summary struct {
// 	[]player.PlayerSummary
// }

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
func (game *Game) AddPlayer(playerType int) {
	switch playerType {
	case c.TYPE_HUMAN:
		game.Table.TakeSeat(player.NewHumanPlayer(), true)
	case c.TYPE_RANDOM:
		game.Table.TakeSeat(player.NewRandomPlayer(), false)
	case c.TYPE_LEARNER:
		game.Table.TakeSeat(player.NewLearnerPlayer(), false)
	case c.TYPE_BASIC:
		game.Table.TakeSeat(player.NewBasicStrategyPlayer(), false)
	case c.TYPE_COUNTER:
		game.Table.TakeSeat(player.NewCardCounterPlayer(), false)
	}
}

// Play is the main game loop
func (game *Game) Play() { // *Summary
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
	// return game.Table.Summarize()
}
