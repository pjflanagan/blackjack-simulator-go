package cards

import (
	"testing"
)

func isEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %+v but got %+v", a, b)
	}
}

func TestIsBlackjack(t *testing.T) {
	type isBlackjackTest []struct {
		caseName string
		hand     *Hand
		result   bool
	}

	blackjackHand := &Hand{
		Cards: []*Card{
			&Card{
				Face: 1,
			},
			&Card{
				Face: 12,
			},
		},
		HasBeenSplit: false,
	}

	splitHand := &Hand{
		Cards: []*Card{
			&Card{
				Face: 1,
			},
			&Card{
				Face: 12,
			},
		},
		HasBeenSplit: true,
	}

	notBlackjackHand := &Hand{
		Cards: []*Card{
			&Card{
				Face: 1,
			},
			&Card{
				Face: 4,
			},
		},
		HasBeenSplit: false,
	}

	threeCardHand := &Hand{
		Cards: []*Card{
			&Card{
				Face: 9,
			},
			&Card{
				Face: 2,
			},

			&Card{
				Face: 10,
			},
		},
		HasBeenSplit: false,
	}

	var testCases = isBlackjackTest{
		{"is blackjack",
			blackjackHand,
			true,
		},
		{"split hand",
			splitHand,
			false,
		},
		{"not blackjack",
			notBlackjackHand,
			false,
		},
		{"three card 21",
			threeCardHand,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			isBlackjack := tc.hand.IsBlackjack()
			isEqual(t, tc.result, isBlackjack)
		})
	}
}
