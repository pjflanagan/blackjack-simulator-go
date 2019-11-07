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
func (table *Table) TakeSeat(newPlayer player.Player, isHuman bool) {
	table.Players = append(table.Players, newPlayer)
	if isHuman {
		table.hasHuman = true
	}
}

// STEP 1: TakeBets --------------------------------------------------------------------------------

// TakeBets goes through all the players and makes them take a bet
// returns true if there is someone playing
func (table *Table) TakeBets() bool {
	fmt.Printf("The count is %d.\n", table.count)
	for _, player := range table.Players {
		if player.CanBet(table.minBet) {
			// if the player can bet then ask them to bet
			player.Bet(table.minBet, table.count)
		} else {
			// if they player can't bet then kick them out
			player.LeaveSeat()
		}
	}
	return table.hasPlayerOfStatus("ANTED")
}

// STEP 2: Deal ------------------------------------------------------------------------------------

// Deal burns a card, makes two passes and gives players and the dealer cards
func (table *Table) Deal() {
	// burn a card
	table.Shoe.Burn()
	for pass := 0; pass < 2; pass++ {
		// make two passes for the deal
		for _, player := range table.Players {
			if player.StatusIs("ANTED", "JEPORADY") {
				// deal each player face up in order
				card := table.takeCard(true)
				player.Deal(0, card)
				if pass == 1 && player.IsBlackjack() {
					// if the player his blackjack on the second pass then reward them
					player.Blackjack()
				}
			}
		}
		// deal the dealer after the players, if it is the first pass keep if face down
		card := table.takeCard(pass == 1)
		table.Dealer.Deal(card)
	}
}

// STEP 3: TakeTurns -------------------------------------------------------------------------------

// TakeTurns makes everyone take turns
func (table *Table) TakeTurns() {
	table.Dealer.PrintHand(table.hasHuman)
	for _, player := range table.Players {
		if player.StatusIs("JEPORADY") {
			// for each player that is playing (in JEPORADY), make them play thier turn
			table.playerTurn(player, 0)
		}
	}
	if table.hasPlayerOfStatus("STAY") {
		// as long as someone is still playing (has STAYed on a hand), the dealer plays
		table.dealerTurn()
	}
}

// playerTurn handles one player's whole turn
func (table *Table) playerTurn(player player.Player, handIdx int) {
	for handIsActive := true; handIsActive; {
		// while the hand is active request the player to move
		move := player.Move(handIdx)
		switch move {
		case "HIT":
			// take a card out and give it to the player, conditional end
			card := table.takeCard(true)
			handIsActive = player.Hit(handIdx, card)
		case "DOUBLE":
			// take a card and double down and end the turn
			card := table.takeCard(true)
			player.DoubleDown(handIdx, card)
			handIsActive = false
		case "SPLIT":
			// have the player split the hand, then take cards and put them in each hand
			player.Split(handIdx)
			card1 := table.takeCard(true)
			card2 := table.takeCard(true)
			player.Hit(handIdx, card1)
			player.Hit(handIdx+1, card2)
			handIsActive = true
		case "STAY":
			// have the player stay, conditional end
			player.Stay(handIdx)
			handIsActive = false
		default:
			// shouldn't happen
		}
	}
	if handIdx < len(player.GetHands())-1 {
		// if there are more hands this player has then move to the next hand
		table.playerTurn(player, handIdx+1)
	}
}

func (table *Table) dealerTurn() {
	// if the dealer needs to take a turn then the dealer shows their card
	revealCard := table.Dealer.RevealCard()
	table.seeCard(revealCard)
	fmt.Printf("Dealer reveals %s. \n", revealCard.Stringify())
	for handIsActive := true; handIsActive; {
		// while the dealer has an active turn
		move := table.Dealer.Move()
		switch move {
		case "HIT":
			// if the dealer said hit then add a card to their hand
			card := table.takeCard(true)
			table.Dealer.Deal(card)
			fmt.Printf("Dealer hits and receives %s.\n", card.Stringify())
			handIsActive = !table.Dealer.DidBust()
		case "STAY":
			// if the dealer said stay then end the turn
			fmt.Printf("Dealer stays.\n")
			handIsActive = false
		default:
			// shouldn't happen
		}
	}
	table.Dealer.PrintHand(table.hasHuman)
}

// STEP 4: Payout ----------------------------------------------------------------------------------

// Payout determines the winnings for each player
func (table *Table) Payout() {
	for _, player := range table.Players {
		if player.StatusIs("STAY") {
			player.Payout(table.Dealer.Hand)
		}
	}
}

// STEP 5: Reset -----------------------------------------------------------------------------------

// Reset resets the table
func (table *Table) Reset() {
	for _, player := range table.Players {
		if !player.StatusIs("OUT") {
			player.Reset(table.minBet)
		}
	}
	if table.Shoe.NeedsShuffle() {
		fmt.Printf("The dealer shuffles the deck.\n")
		table.Shoe.Shuffle()
		table.count = 0
	}
	table.Dealer.Reset()
}

// HELPERS -----------------------------------------------------------------------------------------

// hasPlayerOfStatus returns true when one player's status matches the status
func (table *Table) hasPlayerOfStatus(status string) bool {
	for _, player := range table.Players {
		if player.StatusIs(status) {
			return true
		}
	}
	return false
}

func (table *Table) takeCard(up bool) (card *cards.Card) {
	card = table.Shoe.Take()
	if up {
		// ensure the card is face up then see the card (count the card)
		card.FlipUp()
		table.seeCard(card)
		return
	}
	// otherwise flip it down
	card.FlipDown()
	return
}

func (table *Table) seeCard(card *cards.Card) {
	if card.IsFaceDown() {
		// if the card is face down then don't add it to the count
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
