package player

import (
	"../cards"
	c "../constant"
	"fmt"
	"math/rand"
)

const (
	LEARNER_MAX_PLAYED_HANDS = 500
)

type scenario struct {
	handValue         int // int representing the value of the hand
	handType          int // int representing the type of hand the player had (hard, soft, pair)
	dealerShowingFace int // int representing the face the dealer is showing
	move              int // int representing the move the player made
}

type results struct {
	moveCount map[int]int // map of moves to the number of times they happen
	avgGain   float32
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
	fmt.Printf("%s bets the minumum of %d.\n", player.Name, minBet)
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
		// this would happen if a player gets a 21 after a split (but we shouldn't go to here)
		move = c.MOVE_STAY
	} else {
		move = validMoves[rand.Intn(len(validMoves))]
	}
	player.addScenario(handIdx, dealerHand, move)
	return
}

func (player *LearnerPlayer) addScenario(handIdx int, dealerHand *cards.Hand, move int) {
	// TODO: record to player.handMoves
	handValue, handType := player.Hands[handIdx].Value()
	s := scenario{
		handValue:         handValue,
		handType:          handType,
		dealerShowingFace: dealerHand.ShowingFace(),
		move:              move,
	}
	if player.scenarioResults[s] == nil {
		player.scenarioResults[s] = new(results)
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
	}
}

func (player *LearnerPlayer) addResult(handIdx int, dealerHand *cards.Hand, result int) {
	// TODO: for each move in play hand
	// if it isn't the last move then say continue (not sure if this is best bet for split)
	// if it is the last move then record the result
	handValue, handType := player.Hands[handIdx].Value()
	s := scenario{
		handValue:         handValue,
		handType:          handType,
		dealerShowingFace: dealerHand.ShowingFace(),
		// move:              move,
	}
	if player.scenarioResults[s] == nil {
		player.scenarioResults[s] = new(results)
	}
}
