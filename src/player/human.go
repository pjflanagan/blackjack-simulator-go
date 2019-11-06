package player

import (
	"../cards"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type HumanPlayer struct {
	basePlayer
}

func NewHumanPlayer() *HumanPlayer {
	return &HumanPlayer{
		basePlayer: initBasePlayer(),
	}
}

// STEP 1: Bet -------------------------------------------------------------------------------------

func (player *HumanPlayer) CanBet(minBet int) bool {
	return player.Chips > minBet && player.Active
}

func (player *HumanPlayer) Bet(minBet int, count int) int {
	var bet int
	fmt.Printf("Place bet, you have %d chips [bet 0 to leave table]: ", player.Chips)
	_, _ = fmt.Scanf("%d", &bet)

	if bet == 0 {
		player.LeaveSeat()
		return 0
	} else if bet < 15 {
		fmt.Printf("Bet (%d) is too low sir. ", bet)
		return player.Bet(minBet, count)
	} else if bet > player.Chips {
		fmt.Printf("Bet (%d) is more than what you have sir. ", bet)
		return player.Bet(minBet, count)
	}
	return player.bet(bet)
}

// STEP 3: Move ------------------------------------------------------------------------------------

func (player *HumanPlayer) Move(handIdx int) string {
	player.printHumanHand(handIdx)
	reader := bufio.NewReader(os.Stdin)
	var move string
	// TODO: validMoves
	fmt.Print("Enter [h, s, d, or p]: ")
	move, _ = reader.ReadString('\n')
	move = strings.Replace(move, "\n", "", -1)

	switch move {
	case "h":
		return "HIT"
	case "s":
		return "STAND"
	default:
		fmt.Printf("Move (%s) is invalid pick again.\n", move)
		return player.Move(handIdx)
	}
}

// IsTurnOver returns true when the turn is over
func (player *HumanPlayer) IsTurnOver(handIdx int) bool {
	if player.Hands[handIdx].DidBust() {
		player.printHumanHand(handIdx)
		fmt.Printf("You busted!\n")
		return true
	} else if player.Hands[handIdx].Is21() {
		player.printHumanHand(handIdx)
		fmt.Printf("21! Winner, winner, chicken dinner!\n")
		return true
	}
	return false
}

// STEP 4: Payout ----------------------------------------------------------------------------------

func (player *HumanPlayer) Payout(dealerHand *cards.Hand) {
	for i, hand := range player.Hands {
		result := hand.Result(dealerHand)
		player.payout(i, result)
	}
}

// HELPERS

func (player *HumanPlayer) printHumanHand(handIdx int) {
	fmt.Printf("====== YOUR HAND ======\n")
	fmt.Printf("%s\n", player.HandString(handIdx))
}
