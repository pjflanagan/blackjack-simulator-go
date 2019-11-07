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
		basePlayer: initBasePlayer("You"),
	}
}

// STEP 1: Bet -------------------------------------------------------------------------------------

func (player *HumanPlayer) CanBet(minBet int) bool {
	return player.Chips > minBet && player.Status == "READY"
}

func (player *HumanPlayer) Bet(minBet int, count int) {
	var bet int
	fmt.Printf("Place bet, you have %d chips [bet 0 to leave table]: ", player.Chips)
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
		return "STAY"
	case "d":
		return "DOUBLE"
	case "p":
		return "SPLIT"
	default:
		fmt.Printf("Move (%s) is invalid pick again.\n", move)
		return player.Move(handIdx)
	}
}

// STEP 4: Payout ----------------------------------------------------------------------------------

func (player *HumanPlayer) Payout(dealerHand *cards.Hand) {
	for i, hand := range player.Hands {
		result := hand.Result(dealerHand)
		player.payout(i, result)

		switch result {
		case "BLACKJACK":
			fmt.Printf("You have a blackjack!\n")
		case "WIN":
			fmt.Printf("You won!\n")
		case "PUSH":
			fmt.Printf("You push.\n")
		case "BUST":
			fmt.Printf("You bust.\n")
		case "LOSE":
			fmt.Printf("You lose.\n")
		}
	}
}

// HELPERS

func (player *HumanPlayer) printHumanHand(handIdx int) {
	fmt.Printf("====== YOUR HAND ======\n")
	fmt.Printf("%s\n", player.HandString(handIdx))
}
