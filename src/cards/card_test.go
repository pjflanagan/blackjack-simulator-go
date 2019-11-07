package cards

import (
	"testing"
)

func TestValue(t *testing.T) {
	type cardTest []struct {
		caseName string
		card     *Card
		value    int
	}

	var testCases = cardTest{
		{"ace of spades",
			NewCard(1, 1),
			ACE_VALUE,
		},
		{"ace of hearts",
			NewCard(1, 3),
			ACE_VALUE,
		},
		{"king of clubs",
			NewCard(13, 2),
			FACE_VALUE,
		},
		{"6 of diamonds",
			NewCard(6, 4),
			6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			isEqual(t, tc.value, tc.card.Value())
		})
	}
}

func TestFlip(t *testing.T) {
	card := NewCard(1, 1)

	t.Run("flip", func(t *testing.T) {
		card.FlipUp()
		isEqual(t, false, card.IsFaceDown())
		isEqual(t, true, card.IsFaceUp())
		card.FlipDown()
		isEqual(t, true, card.IsFaceDown())
		isEqual(t, false, card.IsFaceUp())
	})
}
