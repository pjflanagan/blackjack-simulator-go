package cards

import (
	"testing"
)

func TestShoe(t *testing.T) {
	shoe := NewShoe(1)

	type shoeTest []struct {
		caseName         string
		validateFunction func() bool
	}

	var testCases = shoeTest{
		{"burn card",
			func() bool {
				shoe.Burn()
				return len(shoe.deck) == 51 && len(shoe.discards) == 1
			},
		},
		{"take card",
			func() bool {
				shoe.Take()
				return len(shoe.deck) == 50 && len(shoe.discards) == 2
			},
		},
		{"needs shuffle",
			func() bool {
				return shoe.NeedsShuffle() == false
			},
		},
		{"shuffle",
			func() bool {
				shoe.Shuffle()
				return len(shoe.deck) == 52 && len(shoe.discards) == 0
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			isEqual(t, true, tc.validateFunction())
		})
	}
}
