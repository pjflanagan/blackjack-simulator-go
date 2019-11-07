package game

import (
	"../cards"
	"../player"
	"fmt"
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
			if player.GetStatus() == "ANTED" {
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
			if player.GetStatus() == "ANTED" || player.GetStatus() == "JEPORADY" {
				// deal each player in order
				card := table.takeCard(true)
				player.Deal(0, card)
				if pass == 1 && player.IsBlackjack() {
					player.Blackjack()
				}
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
	for _, player := range table.Players {
		if player.GetStatus() == "JEPORADY" {
			table.playerTurn(player, 0)
		}
	}
	if table.hasActivePlayer() {
		table.dealerTurn()
	}
}

func (table *Table) playerTurn(player player.Player, handIdx int) {
	handIsActive := true
	for handIsActive {
		move := player.Move(handIdx)
		handIsActive = table.handleMove(player, handIdx, move)
	}
	if handIdx < len(player.GetHands())-1 {
		table.playerTurn(player, handIdx+1)
	}
}

// handleMove returns true when turn is still active
func (table *Table) handleMove(player player.Player, handIdx int, move string) bool {
	switch move {
	case "HIT":
		// take a card out and give it to the player
		card := table.takeCard(true)
		return player.Hit(handIdx, card)
	case "DOUBLE":
		card := table.takeCard(true)
		player.DoubleDown(handIdx, card)
		return false
	case "SPLIT":
		player.Split(handIdx)
		card1 := table.takeCard(true)
		card2 := table.takeCard(true)
		player.Hit(handIdx, card1)
		player.Hit(handIdx+1, card2)
		return true
	case "STAY":
		player.Stay(handIdx)
		return false
	}
	player.Stay(handIdx)
	return false
}

func (table *Table) hasActivePlayer() bool {
	for _, player := range table.Players {
		if player.GetStatus() == "STAY" {
			return true
		}
	}
	return false
}

func (table *Table) dealerTurn() {
	revealCard := table.Dealer.RevealCard()
	table.seeCard(revealCard)
	fmt.Printf("Dealer reveals %s. \n", revealCard.Stringify())
	activeTurn := true
	for activeTurn {
		move := table.Dealer.Move()
		activeTurn = table.handleDealerMove(move)
	}
	table.Dealer.PrintHand(table.hasHuman)
}

func (table *Table) handleDealerMove(move string) bool {
	switch move {
	case "HIT":
		card := table.takeCard(true)
		table.Dealer.Hand.Add(card)
		fmt.Printf("Dealer hits and receives %s.\n", card.Stringify())
		return !table.Dealer.Hand.DidBust()
	case "STAY":
		fmt.Printf("Dealer stays.\n")
		return false
	}
	return false
}

// STEP 4: Payout ----------------------------------------------------------------------------------

// Payout determines the winnings for each player
func (table *Table) Payout() {
	for _, player := range table.Players {
		if player.GetStatus() == "STAY" {
			player.Payout(table.Dealer.Hand)
		}
	}
}

// STEP 5: Reset -----------------------------------------------------------------------------------

// Reset resets the table
func (table *Table) Reset() {
	for _, player := range table.Players {
		if player.GetStatus() != "OUT" {
			player.Reset(table.minBet)
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
