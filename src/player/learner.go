package player

import (
	"../cards"
	c "../constant"
	"../utils"
	"fmt"
	"os"
	"time"
)

type resultData struct {
	// percentSuccess float32 // map of result to the value change
	successCount int
	count        int
}

func newMoveResultDataMap() map[int]*resultData {
	s := make(map[int]*resultData)
	s[c.MOVE_HIT] = new(resultData)
	s[c.MOVE_STAY] = new(resultData)
	return s
}

// LearnerPlayer extends basePlayer, learner only plays one move per hand, HIT or STAY.
// If busts it records it, otherwise it records the result for the whole hand
type LearnerPlayer struct {
	basePlayer
	shouldRecordHand bool
	originalScenario cards.Scenario                         // the scenario on the deal
	lastMove         int                                    // int represents the last move made
	scenarios        map[cards.Scenario]map[int]*resultData // all of the scenarios played by this player
}

// NewLearnerPlayer returns a new random player with name Random
func NewLearnerPlayer(playerRules *PlayerRules) *LearnerPlayer {
	return &LearnerPlayer{
		basePlayer: initBasePlayer("Learner", playerRules),
		scenarios:  make(map[cards.Scenario]map[int]*resultData),
	}
}

// Bet -------------------------------------------------------------------------------------

// CanBet learners can always bet
func (player *LearnerPlayer) CanBet(minBet int) bool {
	return true
}

// Bet learners bet the minimum
// this bet does not call the parent because we don't want to subtract chips from the learner
func (player *LearnerPlayer) Bet(minBet int, trueCount float32) {
	c.Print("%s bets the minumum of %d chips.\n", player.Name, minBet)
	player.Hands[0].Wager = minBet
	player.Status = c.PLAYER_ANTED
	return
}

// Deal ------------------------------------------------------------------------------------

// CheckDealtHand checks to see if we should record this hand
func (player *LearnerPlayer) CheckDealtHand(dealerHand *cards.Hand, dealerBlackjack bool) {
	player.checkDealtHand(dealerHand, dealerBlackjack)
	if !player.StatusIs(c.PLAYER_JEPORADY) {
		player.shouldRecordHand = false
		if player.StatusIs(c.PLAYER_BLACKJACK) {
			c.Print("%s hit blackjack with a %s!\n", player.Name, player.Hands[0].StringShorthandReadable())
		}
	} else if player.StatusIs(c.PLAYER_JEPORADY) {
		c.Print("%s was dealt %s. (%s, %d)\n",
			player.Name, player.Hands[0].StringShorthandReadable(),
			player.Hands[0].StringScenarioCode(false),
			dealerHand.UpcardValue(),
		)
	}
}

// Move ------------------------------------------------------------------------------------

// Move returns string representing the move
func (player *LearnerPlayer) Move(handIdx int, dealerHand *cards.Hand) (move int) {
	c.Print("%s has %s.\n", player.Name, player.Hands[0].StringSumReadable())

	// only allow the player to move once
	if player.lastMove != 0 {
		return c.MOVE_STAY
	}

	// add the scenario in case it is new
	player.originalScenario, player.shouldRecordHand = player.addScenario(dealerHand)

	// valid moves are only these hits and stay, record times hit bust does not bust and times stay win
	if moves := player.scenarios[player.originalScenario]; moves[c.MOVE_STAY].count > moves[c.MOVE_HIT].count {
		move = c.MOVE_HIT
	} else {
		move = c.MOVE_STAY
	}

	player.lastMove = move
	return
}

func (player *LearnerPlayer) addScenario(dealerHand *cards.Hand) (cards.Scenario, bool) {
	// record this original scenario, if this scenario doesn't exist then add it
	s, shouldRecordHand := cards.NewScenarioFromHands(player.Hands[0], dealerHand, false)
	if player.scenarios[s] == nil {
		// if the scenario is new then add it
		player.scenarios[s] = newMoveResultDataMap()
	}
	return s, shouldRecordHand
}

// Payout ----------------------------------------------------------------------------------

// Payout determines if player should leave and gives the player chips
func (player *LearnerPlayer) Payout(dealerHand *cards.Hand) {
	result := player.Hands[0].Result(dealerHand)
	// add the result to this scenario
	player.addResult(dealerHand, result)
	// call this just to output
	player.resultPayout(0, result)
	player.Chips = 100

	if player.handsPlayed > player.playerRules.MaxHands {
		player.LeaveSeat()
		player.scenariosToCsv() // TODO: move to stats
	}
}

func (player *LearnerPlayer) addResult(dealerHand *cards.Hand, result int) {
	shouldRecordHand := player.shouldRecordHand
	if !shouldRecordHand {
		return
	}
	move := player.lastMove
	s := player.originalScenario

	if move == c.MOVE_HIT {
		// if the move was a hit record data about it
		switch result {
		case c.RESULT_BUST:
			// if they bust then don't add to success, just count
			player.scenarios[s][move].count++
		case c.RESULT_LOSE, c.RESULT_PUSH, c.RESULT_WIN:
			// if the player does not bust, count that as a win for hitting
			// the actual win percent would be then based on the win percent of staying or hitting the following hand
			player.scenarios[s][move].successCount++
			player.scenarios[s][move].count++
		default:
			// print that something went wrong
			c.Print("[ERROR]: hit result was something unexpected. \n")
		}
		// if the scenario was a hit we can also record current scenario if last scenario was a hit, so long as they didn't bust
		s, shouldRecordHand = player.addScenario(dealerHand)
	}

	if !shouldRecordHand {
		// recheck should record hand because we might be checking the final hand
		return
	}
	move = c.MOVE_STAY
	switch result {
	case c.RESULT_WIN:
		player.scenarios[s][move].successCount++
		player.scenarios[s][move].count++
	case c.RESULT_LOSE:
		player.scenarios[s][move].successCount--
		player.scenarios[s][move].count++
	case c.RESULT_PUSH:
		player.scenarios[s][move].count++
	case c.RESULT_BUST:
		// do nothing if they bust, they cannot bust on a stay
	default:
		c.Print("[ERROR]: stay result was something unexpected. \n")
	}
}

// RESET -------------------------------------------------------------------------------------------

func (player *LearnerPlayer) Reset(minBet int) {
	player.Hands = []*cards.Hand{cards.NewHand()}
	player.originalScenario = cards.Scenario{}
	player.shouldRecordHand = false
	player.lastMove = 0
	player.Chips = 100
	player.Status = c.PLAYER_READY
	player.handsPlayed++
}

// SUMMARY -----------------------------------------------------------------------------------------

func (player *LearnerPlayer) scenariosToCsv() {
	// stay and double represent the expected gain, hit represens the odds of not busting
	str := "hand, upcard, stay win, hit doesn't bust, occurances"
	for scenario, result := range player.scenarios {
		if result[c.MOVE_STAY].count+result[c.MOVE_HIT].count == 0 {
			continue
		}
		str = fmt.Sprintf("%s\n%s, %d, %s, %s, %d",
			str,
			scenario.HandString,
			scenario.UpcardValue,
			utils.ToPercent(result[c.MOVE_STAY].successCount, result[c.MOVE_STAY].count),
			utils.ToPercent(result[c.MOVE_HIT].successCount, result[c.MOVE_HIT].count),
			result[c.MOVE_STAY].count+result[c.MOVE_HIT].count,
		)
	}

	f, _ := os.Create(fmt.Sprintf("./out/learn-%d.csv", time.Now().UnixNano()))
	defer f.Close()
	f.WriteString(str)
	f.Sync()
}
