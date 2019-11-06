package game

import (
	"../cards"
	"../player"
)

const ()

// Table represents a blackjack table
type Table struct {
	Shoe     *cards.Shoe
	Players  []player.Player
	Dealer   *Dealer
	minBet   int
	count    int
	hasHuman bool
}

// NewTable returns a table with defaults
func NewTable(minBet int, deckCount int) *Table {
	return &Table{
		Shoe:    cards.NewShoe(deckCount),
		Dealer:  NewDealer(),
		Players: []player.Player{},
		minBet:  minBet,
		count:   0,
	}
}

// STEP 0: Seat ------------------------------------------------------------------------------------

// TakeSeat adds a player to the table
func (table *Table) TakeSeat(newPlayer player.Player) {
	table.Players = append(table.Players, newPlayer)
}

func (table *Table) HasHuman(hasHuman bool) {
	table.hasHuman = hasHuman
}

// STEP 1: TakeBets --------------------------------------------------------------------------------

// TakeBets goes through all the players and makes them take a bet
func (table *Table) TakeBets() (hasActivePlayer bool) {
	for _, player := range table.Players {
		if player.CanBet(table.minBet) {
			player.Bet(table.minBet, table.count)
			if player.IsActive() {
				// if they leave the table while they bet then keep active player false
				hasActivePlayer = true
			}
		} else {
			// if they player can't bet then kick them out
			player.LeaveSeat()
		}
	}
	return
}

// STEP 2: Deal ------------------------------------------------------------------------------------

// Deal burns a card, makes two passes and gives players and the dealer cards
func (table *Table) Deal() {
	// burn a card
	table.Shoe.Burn()
	for pass := 0; pass < 2; pass++ {
		// make two passes for the deal
		for _, player := range table.Players {
			if player.IsActive() {
				// deal each player in order
				card := table.takeCard(true)
				player.Deal(0, card)
			}
		}
		// deal the dealer after the players, if it is the first pass keep if face down
		card := table.takeCard(pass == 1)
		table.Dealer.Hand.Add(card)
	}
}

// STEP 3: TakeTurns -------------------------------------------------------------------------------

// TakeTurns makes everyone take turns
func (table *Table) TakeTurns() {
	table.Dealer.PrintHand(table.hasHuman)
	var hasNonBustPlayer bool
	for _, player := range table.Players {
		if player.IsActive() {
			didBust := table.playerTurn(player)
			if !didBust {
				hasNonBustPlayer = true
			}
		}
	}
	if hasNonBustPlayer {
		table.dealerTurn()
	}
}

func (table *Table) playerTurn(player player.Player) (hasNonBustHand bool) {
	for handIdx := range player.GetHands() {
		activeTurn := true
		for activeTurn {
			move := player.Move(handIdx)
			activeTurn = table.handleMove(player, handIdx, move)
		}
	}
	return
}

func (table *Table) handleMove(player player.Player, handIdx int, move string) bool {
	switch move {
	case "HIT":
		card := table.takeCard(true)
		player.Deal(handIdx, card)
		return !player.IsTurnOver(handIdx)
	case "DOUBLE":
		// player.DoubleDownHand(handIdx)
		card := table.takeCard(true)
		player.Deal(handIdx, card)
		return false
	case "SPLIT":
		return true
	case "STAY":
		return false
	}
	return false
}

func (table *Table) dealerTurn() {
	revealCard := table.Dealer.RevealCard()
	table.seeCard(revealCard)
	activeTurn := true
	for activeTurn {
		move := table.Dealer.Move()
		activeTurn = table.handleDealerMove(move)
	}
}

func (table *Table) handleDealerMove(move string) bool {
	switch move {
	case "HIT":
		card := table.takeCard(true)
		table.Dealer.Hand.Add(card)
		return !table.Dealer.Hand.DidBust()
	case "STAY":
		return false
	}
	return false
}

// STEP 4: Payout ----------------------------------------------------------------------------------

// Payout determines the winnings for each player
func (table *Table) Payout() {
	for _, player := range table.Players {
		if player.IsActive() {
			player.Payout(table.Dealer.Hand)
		}
	}
}

// STEP 5: Reset -----------------------------------------------------------------------------------

// Reset resets the table
func (table *Table) Reset() {
	for _, player := range table.Players {
		if player.IsActive() {
			player.Reset()
		}
	}
	if table.Shoe.NeedsShuffle() {
		table.Shoe.Shuffle()
	}
	table.Dealer.Reset()
}

// HELPERS -----------------------------------------------------------------------------------------

func (table *Table) takeCard(up bool) (card *cards.Card) {
	card = table.Shoe.Take()
	if up {
		table.seeCard(card)
	} else {
		card.FlipDown()
	}
	return
}

func (table *Table) seeCard(card *cards.Card) {
	if card.FaceDown {
		return
	}
	value := card.Value()
	switch value {
	case cards.ACE_VALUE, 10:
		table.count--
	case 2, 3, 4, 5, 6:
		table.count++
	case 7, 8, 9:
	}
}
