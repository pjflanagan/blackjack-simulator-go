package game

import (
	c "../constant"
	"../player"
	"../stats"
)

const ()

// Game represents a blackjack game
type Game struct {
	Table *Table
}

// NewGame returns a game with defaults
func NewGame(minBet int, deckCount int, runNumber int) *Game {
	return &Game{
		Table: NewTable(minBet, deckCount, runNumber),
	}
}

// AddPlayer adds a new player to the game
func (game *Game) AddPlayer(playerType int, playerRules *player.PlayerRules) {
	switch playerType {
	case c.TYPE_HUMAN:
		game.Table.TakeSeat(player.NewHumanPlayer(playerRules))
	case c.TYPE_RANDOM:
		game.Table.TakeSeat(player.NewRandomPlayer(playerRules))
	case c.TYPE_LEARNER:
		game.Table.TakeSeat(player.NewLearnerPlayer(playerRules))
	case c.TYPE_BASIC:
		game.Table.TakeSeat(player.NewBasicStrategyPlayer(playerRules))
	case c.TYPE_COUNTER:
		game.Table.TakeSeat(player.NewCardCounterPlayer(playerRules))
	}
}

// Play is the main game loop
func (game *Game) Play() []*stats.Stats {
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
	return game.Table.Summarize()
}
