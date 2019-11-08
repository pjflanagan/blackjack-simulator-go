package game

import (
	"../cards"
	c "../constant"
	"../player"
	"fmt"
)

const ()

// Table represents a blackjack table
type Table struct {
	Shoe      *cards.Shoe
	Players   []player.Player
	Dealer    *Dealer
	minBet    int
	count     int
	hasHuman  bool
	handCount int
}

// NewTable returns a table with defaults
func NewTable(minBet int, deckCount int) *Table {
	return &Table{
		Shoe:      cards.NewShoe(deckCount),
		Dealer:    NewDealer(),
		Players:   []player.Player{},
		minBet:    minBet,
		count:     0,
		handCount: 0,
	}
}

// Seat ------------------------------------------------------------------------------------

// TakeSeat adds a player to the table
func (table *Table) TakeSeat(newPlayer player.Player, isHuman bool) {
	table.Players = append(table.Players, newPlayer)
	if isHuman {
		table.hasHuman = true
	}
}

// TakeBets --------------------------------------------------------------------------------

// TakeBets goes through all the players and makes them take a bet
// returns true if there is someone playing
func (table *Table) TakeBets() bool {
	if !table.hasPlayerOfStatus(c.PLAYER_READY) {
		return false
	}
	fmt.Printf("\n\n======= HAND %d =======\n", table.handCount)
	fmt.Printf("The count is %d.\n", table.count)
	fmt.Printf("\n == Bet ==\n")
	for _, player := range table.Players {
		if player.CanBet(table.minBet) {
			// if the player can bet then ask them to bet
			player.Bet(table.minBet, table.count)
		} else {
			// if they player can't bet then kick them out
			player.LeaveSeat()
		}
	}
	return table.hasPlayerOfStatus(c.PLAYER_ANTED)
}

// Deal ------------------------------------------------------------------------------------

// Deal burns a card, makes two passes and gives players and the dealer cards
func (table *Table) Deal() {
	fmt.Printf("\n == Deal ==\n")
	// burn a card
	table.Shoe.Burn()
	for pass := 0; pass < 2; pass++ {
		// make two passes for the deal
		for _, player := range table.Players {
			if player.StatusIs(c.PLAYER_ANTED, c.PLAYER_JEPORADY) {
				// deal each player face up in order
				card := table.takeCard(true)
				player.Deal(0, card)
				if pass == 1 {
					// check if the player got blackjack and output what the hand was
					player.CheckDealtHand()
				}
			}
		}
		// deal the dealer after the players, if it is the first pass flip it face up, second pass is face down
		card := table.takeCard(pass == 0)
		table.Dealer.Deal(card)
	}
	table.Dealer.PrintHand(table.hasHuman)
}

// TakeTurns -------------------------------------------------------------------------------

// TakeTurns makes everyone take turns
func (table *Table) TakeTurns() {
	fmt.Printf("\n == Turns ==\n")
	for _, player := range table.Players {
		if player.StatusIs(c.PLAYER_JEPORADY) {
			// for each player that is playing (in JEPORADY), make them play thier turn
			table.playerTurn(player, 0)
		}
	}
	if table.hasPlayerOfStatus(c.PLAYER_STAY) {
		// as long as someone is still playing (has STAYed on a hand), the dealer plays
		table.dealerTurn()
	}
}

// playerTurn handles one player's whole turn
func (table *Table) playerTurn(player player.Player, handIdx int) {
	for handIsActive := true; handIsActive; {
		// while the hand is active request the player to move
		move := player.Move(handIdx, table.Dealer.Hand)
		switch move {
		case c.MOVE_HIT:
			// take a card out and give it to the player, conditional end
			card := table.takeCard(true)
			handIsActive = player.Hit(handIdx, card)
		case c.MOVE_DOUBLE:
			// take a card and double down and end the turn
			card := table.takeCard(true)
			player.DoubleDown(handIdx, card)
			handIsActive = false
		case c.MOVE_SPLIT:
			// have the player split the hand, then take cards and put them in each hand
			player.Split(handIdx)
			card1 := table.takeCard(true)
			card2 := table.takeCard(true)
			player.Hit(handIdx, card1)
			player.Hit(handIdx+1, card2)
			handIsActive = true
		case c.MOVE_STAY:
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
	for handIsActive := true; handIsActive; {
		// while the dealer has an active turn
		move := table.Dealer.Move()
		switch move {
		case c.MOVE_HIT:
			// if the dealer said hit then add a card to their hand
			card := table.takeCard(true)
			table.Dealer.Hit(card)
			handIsActive = !table.Dealer.DidBust()
		case c.MOVE_STAY:
			// if the dealer said stay then end the turn
			table.Dealer.Stay()
			handIsActive = false
		default:
			// shouldn't happen
		}
	}
	table.Dealer.PrintHand(table.hasHuman)
}

// Payout ----------------------------------------------------------------------------------

// Payout determines the winnings for each player
func (table *Table) Payout() {
	fmt.Printf("\n == Payout ==\n")
	for _, player := range table.Players {
		if !player.StatusIs(c.PLAYER_OUT) {
			player.Payout(table.Dealer.Hand)
		}
	}
}

// Reset -----------------------------------------------------------------------------------

// Reset resets the table
func (table *Table) Reset() {
	table.handCount++
	for _, player := range table.Players {
		if !player.StatusIs(c.PLAYER_OUT) {
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
func (table *Table) hasPlayerOfStatus(status int) bool {
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
	case cards.ACE_VALUE, cards.FACE_VALUE:
		table.count--
	case 2, 3, 4, 5, 6:
		table.count++
	case 7, 8, 9:
	}
}
