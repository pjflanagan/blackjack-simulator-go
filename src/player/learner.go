package player

import (
	"../cards"
	c "../constant"
	"../utils"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	LEARNER_MAX_PLAYED_HANDS = 2000
)

// TODO: these should maybe be in the card class
type scenario struct {
	handString         string // string representing the player's hand
	dealerShowingValue int    // value of the dealer's upcard
	move               int    // int representing the move the player made
}

type results struct {
	resultCount map[int]int // map of result to the number of times they happen
}

func newResults() *results {
	return &results{
		resultCount: make(map[int]int),
	}
}

// LearnerPlayer extends basePlayer
type LearnerPlayer struct {
	basePlayer
	scenarioResults map[scenario]*results
	handMoves       [][]int // int represents the moves made for each hand
	playedHands     int
}

// NewLearnerPlayer returns a new random player with name Random
func NewLearnerPlayer() *LearnerPlayer {
	learnerPlayer := &LearnerPlayer{
		basePlayer:      initBasePlayer("Learner"),
		scenarioResults: make(map[scenario]*results),
	}
	learnerPlayer.basePlayer.child = learnerPlayer
	return learnerPlayer
}

// Bet -------------------------------------------------------------------------------------

// CanBet learners can always bet
func (player *LearnerPlayer) CanBet(minBet int) bool {
	return true
}

// Bet learners bet the minimum
// this bet does not call the parent because we don't want to subtract chips from the learner
func (player *LearnerPlayer) Bet(minBet int, count int) {
	fmt.Printf("%s bets the minumum of %d chips.\n", player.Name, minBet)
	player.playedHands++
	player.Hands[0].Wager = minBet
	player.Status = c.PLAYER_ANTED
	return
}

// Move ------------------------------------------------------------------------------------

// Move returns string representing the move
func (player *LearnerPlayer) Move(handIdx int, dealerHand *cards.Hand) (move int) {
	fmt.Printf("%s has %s.\n", player.Name, player.Hands[handIdx].ShorthandSumString())
	validMoves := player.Hands[handIdx].GetValidMoves(100)
	if len(validMoves) == 0 {
		// this would happen if a player gets a 21 after a split
		move = c.MOVE_STAY
	} else if utils.Contains(validMoves, c.MOVE_SPLIT) {
		move = c.MOVE_SPLIT
	} else {
		move = validMoves[rand.Intn(len(validMoves))]
	}
	player.addScenario(handIdx, dealerHand, move)
	return
}

func (player *LearnerPlayer) addScenario(handIdx int, dealerHand *cards.Hand, move int) {
	scenarioString := player.Hands[handIdx].ScenarioString()
	if scenarioString == "" {
		return
	}
	s := scenario{
		handString:         scenarioString,
		dealerShowingValue: dealerHand.ShowingValue(),
		move:               move,
	}
	if player.scenarioResults[s] == nil {
		player.scenarioResults[s] = newResults()
	}
	if len(player.handMoves) == handIdx {
		player.handMoves = append(player.handMoves, []int{move})
	} else {
		player.handMoves[handIdx] = append(player.handMoves[handIdx], move)
	}
}

// Payout ----------------------------------------------------------------------------------

// Payout determines if player should leave and gives the player chips
func (player *LearnerPlayer) Payout(dealerHand *cards.Hand) {
	for i, hand := range player.Hands {
		result := hand.Result(dealerHand)
		player.addResult(i, dealerHand, result)
		player.resultPayout(i, result)
	}

	if player.playedHands > LEARNER_MAX_PLAYED_HANDS {
		player.LeaveSeat()
		player.Summarize()
	}

	// fmt.Printf("\n%+v\n", player.scenarioResultsToCsv())
}

func (player *LearnerPlayer) addResult(handIdx int, dealerHand *cards.Hand, result int) {
	scenarioString := player.Hands[handIdx].ScenarioString()
	if scenarioString == "" {
		return
	}
	for _, move := range player.handMoves[handIdx] {
		// for each move they made this hand
		s := scenario{
			handString:         scenarioString,
			dealerShowingValue: dealerHand.ShowingValue(),
			move:               move,
		}
		if player.scenarioResults[s] == nil {
			// if this scenario doesn't exist then add it
			player.scenarioResults[s] = newResults()
		}
		player.scenarioResults[s].resultCount[result]++
	}
}

func (player *LearnerPlayer) Summarize() (str string) {
	player.scenarioResultsToCsv()
	return ""
}

var moveStringMap = map[int]string{
	c.MOVE_SPLIT:  "split",
	c.MOVE_STAY:   "stay",
	c.MOVE_HIT:    "hit",
	c.MOVE_DOUBLE: "double",
}

func (player *LearnerPlayer) scenarioResultsToCsv() {
	str := "hand,upcard,move,win,lose,push"
	for scenario, result := range player.scenarioResults {
		str = fmt.Sprintf("%s\n%s,%d,%s,%d,%d,%d",
			str,
			scenario.handString,
			scenario.dealerShowingValue,
			moveStringMap[scenario.move],
			result.resultCount[c.RESULT_WIN],
			result.resultCount[c.RESULT_BUST]+result.resultCount[c.RESULT_LOSE],
			result.resultCount[c.RESULT_PUSH],
		)
	}

	f, _ := os.Create(fmt.Sprintf("./out/learn-%d.csv", time.Now().UnixNano()))
	defer f.Close()
	f.WriteString(str)
	f.Sync()
}
