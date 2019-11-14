package player

import (
	"../cards"
	c "../constant"
	"../utils"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// HumanPlayer extends basePlayer
type HumanPlayer struct {
	basePlayer
}

// NewHumanPlayer returns a new human player with name You
func NewHumanPlayer(playerRules *PlayerRules) *HumanPlayer {
	return &HumanPlayer{
		basePlayer: initBasePlayer("Human", playerRules),
	}
}

// Bet -------------------------------------------------------------------------------------

// Bet prompts a player to bet
func (player *HumanPlayer) Bet(minBet int, trueCount float32) {
	var bet int
	c.Print("%s have %d chips, place bet or 0 to leave: ", player.Name, player.Chips)
	_, _ = fmt.Scanf("%d", &bet)

	if bet == 0 {
		player.LeaveSeat()
		return
	} else if bet < minBet {
		c.Print("Bet (%d) is too low sir. ", bet)
		player.Bet(minBet, trueCount)
		return
	} else if bet > player.Chips {
		c.Print("Bet (%d) is more than what you have sir. ", bet)
		player.Bet(minBet, trueCount)
		return
	}
	player.bet(bet)
	return
}

// Move ------------------------------------------------------------------------------------

// Move returns string representing the move
func (player *HumanPlayer) Move(handIdx int, dealerHand *cards.Hand) (move int) {
	// player.PrintVisualHand(handIdx)
	validMoves := player.Hands[handIdx].GetValidMoves(player.Chips)
	if len(validMoves) == 0 {
		// this would happen if a player gets a 21 after a split (but we shouldn't go to here)
		return c.MOVE_STAY
	}

	reader := bufio.NewReader(os.Stdin)
	var input string
	c.Print("Human has %s. Enter move %s: ", player.Hands[handIdx].StringSumReadable(), getValidMovePrompt(validMoves))
	input, _ = reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)

	switch input {
	case "h":
		move = c.MOVE_HIT
	case "s":
		move = c.MOVE_STAY
	case "d":
		move = c.MOVE_DOUBLE
	case "p":
		move = c.MOVE_SPLIT
	default:
		c.Print("Move (%s) is invalid pick again.\n", input)
		move = player.Move(handIdx, dealerHand)
	}

	if !utils.Contains(validMoves, move) {
		c.Print("Move (%s) is invalid pick again.\n", input)
		move = player.Move(handIdx, dealerHand)
	}
	return
}

func getValidMovePrompt(validMoves []int) string {
	var validMoveChars []string
	for _, move := range validMoves {
		switch move {
		case c.MOVE_HIT:
			validMoveChars = append(validMoveChars, "H")
		case c.MOVE_STAY:
			validMoveChars = append(validMoveChars, "S")
		case c.MOVE_DOUBLE:
			validMoveChars = append(validMoveChars, "D")
		case c.MOVE_SPLIT:
			validMoveChars = append(validMoveChars, "P")
		}
	}
	return fmt.Sprintf("%s", validMoveChars)
}

// Payout ----------------------------------------------------------------------------------

// Payout print's message hand handles the payout
func (player *HumanPlayer) Payout(dealerHand *cards.Hand) {
	for i, hand := range player.Hands {
		result := hand.Result(dealerHand)
		player.resultPayout(i, result)
	}
}
