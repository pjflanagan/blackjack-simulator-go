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
func NewHumanPlayer() *HumanPlayer {
	humanPlayer := &HumanPlayer{
		basePlayer: initBasePlayer("You"),
	}
	humanPlayer.basePlayer.child = humanPlayer
	return humanPlayer
}

// STEP 1: Bet -------------------------------------------------------------------------------------

// CanBet returns true when a player can bet
func (player *HumanPlayer) CanBet(minBet int) bool {
	return player.Chips > minBet && player.Status == c.PLAYER_READY
}

// Bet prompts a player to bet
func (player *HumanPlayer) Bet(minBet int, count int) {
	var bet int
	fmt.Printf("You have %d chips, place bet or 0 to leave: ", player.Chips)
	_, _ = fmt.Scanf("%d", &bet)

	if bet == 0 {
		player.LeaveSeat()
		return
	} else if bet < 15 {
		fmt.Printf("Bet (%d) is too low sir. ", bet)
		player.Bet(minBet, count)
		return
	} else if bet > player.Chips {
		fmt.Printf("Bet (%d) is more than what you have sir. ", bet)
		player.Bet(minBet, count)
		return
	}
	player.bet(0, bet)
	return
}

// Blackjack handles when a player hits blackjack
func (player *HumanPlayer) Blackjack() {
	fmt.Printf("%s hit blackjack!\n", player.Name)
	player.printHumanHand(0)
	player.blackjack()
}

// STEP 3: Move ------------------------------------------------------------------------------------

// Move returns string representing the move
func (player *HumanPlayer) Move(handIdx int) (move int) {
	player.printHumanHand(handIdx)
	validMoves := player.Hands[handIdx].GetValidMoves(player.Chips)

	reader := bufio.NewReader(os.Stdin)
	var input string
	fmt.Printf("Enter %s: ", getValidMovePrompt(validMoves))
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
		fmt.Printf("Move (%s) is invalid pick again.\n", input)
		move = player.Move(handIdx)
	}

	if !utils.Contains(validMoves, move) {
		fmt.Printf("Move (%s) is invalid pick again.\n", input)
		move = player.Move(handIdx)
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

// Hit returns true if hand is still active
func (player *HumanPlayer) Hit(handIdx int, card *cards.Card) bool {
	fmt.Printf("%s received %s.\n", player.Name, card.Stringify())
	return player.hit(handIdx, card)
}

// Split splits the player's hand
func (player *HumanPlayer) Split(handIdx int) {
	fmt.Printf("%s split.\n", player.Name)
	player.split(handIdx)
}

// DoubleDown doubles down
func (player *HumanPlayer) DoubleDown(handIdx int, card *cards.Card) {
	fmt.Printf("%s double down and receive %s.\n", player.Name, card.Stringify())
	player.doubleDown(handIdx, card)
}

// Bust busts the players hand and sets the status
func (player *HumanPlayer) Bust(handIdx int) {
	fmt.Printf("%s bust and lose %d.\n", player.Name, player.Hands[handIdx].Wager)
	player.bust(handIdx)
}

// Stay returns true if the player's turn is still active
func (player *HumanPlayer) Stay(handIdx int) {
	fmt.Printf("%s stay.\n", player.Name)
	player.stay(handIdx)
}

// STEP 4: Payout ----------------------------------------------------------------------------------

// Payout print's message hand handles the payout
func (player *HumanPlayer) Payout(dealerHand *cards.Hand) {
	for i, hand := range player.Hands {
		result := hand.Result(dealerHand)

		switch result {
		case c.RESULT_BLACKJACK:
			// do not call payout for blackjack, money has already been given
			fmt.Printf("You had a blackjack!\n")
		case c.RESULT_WIN:
			fmt.Printf("You won!\n")
			player.payout(i, result)
		case c.RESULT_PUSH:
			fmt.Printf("You push.\n")
			player.payout(i, result)
		case c.RESULT_BUST:
			// do not call payout for bust, money has already been taken
			fmt.Printf("You busted.\n")
		case c.RESULT_LOSE:
			fmt.Printf("You lose.\n")
			player.payout(i, result)
		}
	}
}

// HELPERS

func (player *HumanPlayer) printHumanHand(handIdx int) {
	fmt.Printf("\n====== YOUR HAND ======\n")
	fmt.Printf("You have a %s.\n", player.Hands[handIdx].ShorthandSumString())
	fmt.Printf("%s\n", player.Hands[handIdx].LongformString())
}
