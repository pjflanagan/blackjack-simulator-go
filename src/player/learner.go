package player

import (
	"../cards"
	c "../constant"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	LEARNER_MAX_PLAYED_HANDS = 10000
)

type results struct {
	moveWinPercent map[int]float32 // map of result to the value change
	moveCount      map[int]int
}

func newResults(move int) *results {
	return &results{
		moveWinPercent: map[int]float32{
			move: float32(0),
		},
		moveCount: map[int]int{
			move: 0,
		},
	}
}

// LearnerPlayer extends basePlayer, learner only plays one move per hand, HIT, DOUBLE, STAY (it will never split).
// Every time it gets to an existing scenario it does whichever one was done last.
// If busts it records it, otherwise it records the result for the whole hand

type LearnerPlayer struct {
	basePlayer
	originalScenario cards.Scenario              // the scenario on the deal
	lastMove         int                         // int represents the last move made
	scenarioResults  map[cards.Scenario]*results // all of the scenarios played by this player
	playedHands      int                         // count of all the hands played so we quit eventually
}

// NewLearnerPlayer returns a new random player with name Random
func NewLearnerPlayer() *LearnerPlayer {
	return &LearnerPlayer{
		basePlayer:      initBasePlayer("Learner"),
		scenarioResults: make(map[cards.Scenario]*results),
	}
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

	// only allow the player to move once
	if player.lastMove != 0 {
		return c.MOVE_STAY
	}

	// valid moves are only these three
	validMoves := []int{c.MOVE_HIT, c.MOVE_DOUBLE, c.MOVE_STAY}
	move = validMoves[rand.Intn(len(validMoves))]
	player.lastMove = move

	// record this original scenario, if this scenario doesn't exist then add it
	s, _ := cards.NewScenario(player.Hands[handIdx], dealerHand)
	if player.scenarioResults[s] == nil {
		player.scenarioResults[s] = newResults(move)
	}
	player.originalScenario = s

	return
}

// Payout ----------------------------------------------------------------------------------

// Payout determines if player should leave and gives the player chips
func (player *LearnerPlayer) Payout(dealerHand *cards.Hand) {
	for i, hand := range player.Hands {
		result := hand.Result(dealerHand)
		// add the result to this scenario
		player.addResult(i, dealerHand, result)
		// call this just to output
		player.resultPayout(i, result)
		player.Chips = 100
	}

	if player.playedHands > LEARNER_MAX_PLAYED_HANDS {
		player.LeaveSeat()
		player.Summarize()
	}
}

func (player *LearnerPlayer) addResult(handIdx int, dealerHand *cards.Hand, result int) {
	move := player.lastMove
	s := player.originalScenario

	// record what happens for the move they made for the scenario they were in
	if move == c.MOVE_HIT {
		switch result {
		case c.RESULT_BUST:
			reAverage(player.scenarioResults[s], move, -1)
		default:
			// if the player does not bust, count that as a win for hitting
			// the actual win percent would be then based on the win percent of staying or hitting the following hand
			reAverage(player.scenarioResults[s], move, 1)
		}
	}
	switch result {
	case c.RESULT_WIN:
		reAverage(player.scenarioResults[s], move, 1)
	case c.RESULT_BUST, c.RESULT_LOSE:
		reAverage(player.scenarioResults[s], move, -1)
	case c.RESULT_PUSH:
		reAverage(player.scenarioResults[s], move, 0)
	}
}

func reAverage(scenarioResults *results, move int, addToAverage int) {
	total := scenarioResults.moveWinPercent[move]*float32(scenarioResults.moveCount[move]) + float32(addToAverage)
	scenarioResults.moveCount[move]++
	scenarioResults.moveWinPercent[move] = total / float32(scenarioResults.moveCount[move])
}

// RESET -------------------------------------------------------------------------------------------

func (player *LearnerPlayer) Reset(minBet int) {
	player.Hands = []*cards.Hand{cards.NewHand()}
	player.originalScenario = cards.Scenario{}
	player.lastMove = 0
	player.Chips = 100
	player.Status = c.PLAYER_READY
}

// SUMMARY -----------------------------------------------------------------------------------------

func (player *LearnerPlayer) Summarize() (str string) {
	player.scenarioResultsToCsv()
	return ""
}

func (player *LearnerPlayer) scenarioResultsToCsv() {
	// stay and double represent the expected gain, hit represens the odds of not busting
	str := "hand, upcard, stay, double, hit, occurances"
	for scenario, result := range player.scenarioResults {
		var total int
		for _, amount := range result.moveCount {
			total += amount
		}
		str = fmt.Sprintf("%s\n%s, %d, %s, %s, %s, %d",
			str,
			scenario.HandString,
			scenario.UpcardValue,
			toPercent(result.moveWinPercent[c.MOVE_STAY]),
			toPercent(result.moveWinPercent[c.MOVE_DOUBLE]*2),
			toPercent(result.moveWinPercent[c.MOVE_HIT]),
			total,
		)
	}

	f, _ := os.Create(fmt.Sprintf("./out/learn-%d.csv", time.Now().UnixNano()))
	defer f.Close()
	f.WriteString(str)
	f.Sync()
}

func toPercent(num float32) string {
	return fmt.Sprintf("%.2f%%", num*100)
}
