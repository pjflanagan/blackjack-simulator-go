package constant

import (
	"fmt"
)

const (
	TYPE_HUMAN = iota
	TYPE_RANDOM
	TYPE_LEARNER
	TYPE_BASIC
	TYPE_COUNTER

	// output modes
	OUTPUT_HUMAN
	OUTPUT_LOG
	OUTPUT_NONE

	// player status
	PLAYER_READY
	PLAYER_ANTED
	PLAYER_JEPORADY
	PLAYER_BLACKJACK
	PLAYER_BUST
	PLAYER_STAY
	PLAYER_OUT

	// hand status
	HAND_SOFT
	HAND_HARD
	HAND_PAIR

	// result status
	RESULT_BLACKJACK
	RESULT_BUST
	RESULT_STAY
	RESULT_PUSH
	RESULT_WIN
	RESULT_LOSE

	// move status
	MOVE_HIT
	MOVE_STAY
	MOVE_DOUBLE
	MOVE_SPLIT

	// default
	DEFAULT_MIN   = 10
	DEFAULT_DECKS = 6
	DEFAULT_CHIPS = 100
)

var OUTPUT_MODE int

func SetOutputMode(mode int) {
	OUTPUT_MODE = mode
}

func Print(base string, inserts ...interface{}) {
	if OUTPUT_MODE != OUTPUT_NONE {
		fmt.Printf(base, inserts...)
	}
}
