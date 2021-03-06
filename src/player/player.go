package player

import (
	"../cards"
	c "../constant"
	"../stats"
)

// Player is the base class for all players (excluding dealer)
type Player interface {
	// Bet
	CanBet(minBet int) bool
	Bet(minBet int, trueCount float32)
	// Deal
	Deal(handIdx int, card *cards.Card)
	CheckDealtHand(dealerHand *cards.Hand, dealerBlackjack bool)
	// Move
	Move(handIdx int, dealerHand *cards.Hand) int
	Hit(handIdx int, card *cards.Card) bool
	SplitHit(handIdx int, card *cards.Card) bool
	Split(handIdx int)
	DoubleDown(handIdx int, card *cards.Card)
	Stay(handIdx int)
	Bust(handIdx int)
	// Payout
	Payout(dealerHand *cards.Hand)
	// Reset
	Reset(minBet int)
	LeaveSeat()
	// Helpers
	GetHands() []*cards.Hand
	GetHand(handIdx int) *cards.Hand
	GetChips() int
	StatusIs(statuses ...int) bool
	PrintVisualHand(handIdx int)
	Summarize() *stats.Stats
}

type PlayerRules struct {
	StartingChips int // chips they start with
	LeavingChips  int // chips they'll leave at
	MaxHands      int // hand they'll leave at, -1 for infinite
}

type basePlayer struct {
	Name        string
	Hands       []*cards.Hand
	Chips       int
	Status      int
	handsPlayed int
	playerRules *PlayerRules
}

func initBasePlayer(name string, playerRules *PlayerRules) basePlayer {
	if playerRules.StartingChips == 0 {
		playerRules.StartingChips = c.DEFAULT_CHIPS
	}
	return basePlayer{
		Name:        name,
		Hands:       []*cards.Hand{cards.NewHand()},
		Chips:       playerRules.StartingChips,
		Status:      c.PLAYER_READY,
		playerRules: playerRules,
	}
}

// Bet -------------------------------------------------------------------------------------

func (player *basePlayer) CanBet(minBet int) bool {
	return player.Chips >= minBet && player.StatusIs(c.PLAYER_READY)
}

func (player *basePlayer) Bet(minBet int, trueCount float32) {
	bet := minBet
	c.Print("%s bets the minimum %d of %d chips available.\n", player.Name, bet, player.Chips)
	player.bet(bet)
	return
}

// bet is the initial bet always on hand 0
func (player *basePlayer) bet(bet int) {
	if player.StatusIs(c.PLAYER_OUT) {
		return
	}
	player.Chips -= bet
	player.Hands[0].Wager = bet
	player.Status = c.PLAYER_ANTED
}

// Deal ------------------------------------------------------------------------------------

// Deal adds a card to the player's hand
func (player *basePlayer) Deal(handIdx int, card *cards.Card) {
	player.Hands[handIdx].Add(card)
	player.Status = c.PLAYER_JEPORADY
}

// CheckDealtHand prints a statment with what they we're dealt
func (player *basePlayer) CheckDealtHand(dealerHand *cards.Hand, dealerBlackjack bool) {
	player.checkDealtHand(dealerHand, dealerBlackjack)
	if player.StatusIs(c.PLAYER_BLACKJACK) {
		c.Print("%s hit blackjack with a %s!\n", player.Name, player.Hands[0].StringShorthandReadable())
	} else if player.StatusIs(c.PLAYER_JEPORADY) {
		c.Print("%s was dealt %s.\n", player.Name, player.Hands[0].StringShorthandReadable())
	}
}

func (player *basePlayer) checkDealtHand(dealerHand *cards.Hand, dealerBlackjack bool) {
	if player.Hands[0].IsBlackjack() {
		if dealerBlackjack {
			// if they both hit blackjack then push
			player.Status = c.PLAYER_BLACKJACK
		} else {
			// if just player hit blackjack then pay them
			player.blackjack()
		}
	} else {
		if dealerBlackjack {
			// if the dealer hits blackjack then lose
			player.Status = c.PLAYER_BUST
		} else {
			// otherwise we're playing the hand
			player.Status = c.PLAYER_JEPORADY
		}
	}
}

func (player *basePlayer) blackjack() {
	player.Chips += player.payout(0, c.RESULT_BLACKJACK) // give the money right away
	player.Status = c.PLAYER_BLACKJACK
}

// Turn ------------------------------------------------------------------------------------

func (player *basePlayer) validMoves() []string {
	return []string{}
}

// Hit returns true if hand is still active
func (player *basePlayer) Hit(handIdx int, card *cards.Card) bool {
	c.Print("%s hits and receives %s.\n", player.Name, card.StringShorthand())
	return player.hit(handIdx, card)
}

// SplitHit returns true if hand is still active
func (player *basePlayer) SplitHit(handIdx int, card *cards.Card) bool {
	c.Print("%s receives %s.\n", player.Name, card.StringShorthand())
	return player.hit(handIdx, card)
}

// Hit returns true when hand is still active
func (player *basePlayer) hit(handIdx int, card *cards.Card) bool {
	player.Deal(handIdx, card)
	if player.Hands[handIdx].DidBust() {
		// if they bust then determine if turn is really over
		player.Bust(handIdx)
		return false
	} else if player.Hands[handIdx].Is21() {
		// if they hit 21 then this hand is over
		player.stay(handIdx)
		return false
	}
	// turn is still active, status is JEPORADY
	return true
}

// Split splits the player's hand
func (player *basePlayer) Split(handIdx int) {
	c.Print("%s splits.\n", player.Name)
	player.split(handIdx)
}

func (player *basePlayer) split(handIdx int) {
	player.Chips -= player.Hands[handIdx].Wager
	splitHand := player.Hands[handIdx].Split()

	if handIdx == len(player.Hands)-1 {
		// if at the end of the array append the new hand to the end
		player.Hands = append(player.Hands, splitHand)
	} else {
		// if not at the end of the array then do something shifty
		// make space in the array for a new element
		player.Hands = append(player.Hands, nil)
		// copy over elements sourced from handIdx to one over
		copy(player.Hands[handIdx+2:], player.Hands[handIdx+1:])
		player.Hands[handIdx+1] = splitHand
	}

	player.Status = c.PLAYER_JEPORADY
}

// DoubleDown doubles down
func (player *basePlayer) DoubleDown(handIdx int, card *cards.Card) {
	c.Print("%s doubles down and receives %s.\n", player.Name, card.StringShorthand())
	player.doubleDown(handIdx, card)
}

func (player *basePlayer) doubleDown(handIdx int, card *cards.Card) {
	player.Hands[handIdx].Add(card)
	player.Chips -= player.Hands[handIdx].Wager
	player.Hands[handIdx].Wager *= 2
	if player.Hands[handIdx].DidBust() {
		player.Bust(handIdx)
	} else {
		player.stay(handIdx)
	}
}

// Stay returns true if the player's turn is still active
func (player *basePlayer) Stay(handIdx int) {
	c.Print("%s stays.\n", player.Name)
	player.stay(handIdx)
}

// Returns true if the player's turn is still active
func (player *basePlayer) stay(handIdx int) {
	if handIdx == len(player.Hands)-1 {
		player.Status = c.PLAYER_STAY
	} else {
		player.Status = c.PLAYER_JEPORADY
	}
}

// Bust busts the players hand and sets the status
func (player *basePlayer) Bust(handIdx int) {
	c.Print("%s busts and loses %d chips.\n", player.Name, player.Hands[handIdx].Wager)
	player.bust(handIdx)
}

// Returns true if the player's hand is still active
func (player *basePlayer) bust(handIdx int) {
	// no need to do payout they wont recieve money for this
	if handIdx == len(player.Hands)-1 {
		player.Status = c.PLAYER_BUST
	} else {
		player.Status = c.PLAYER_JEPORADY
	}
}

// Payout ----------------------------------------------------------------------------------

// Payout print's message hand handles the payout
func (player *basePlayer) Payout(dealerHand *cards.Hand) {
	for i, hand := range player.Hands {
		result := hand.Result(dealerHand)
		player.resultPayout(i, result)
	}

	if player.Chips > player.playerRules.LeavingChips || (player.playerRules.MaxHands != 0 && player.handsPlayed > player.playerRules.MaxHands) {
		player.LeaveSeat()
	}
}

// payout does the math for the payout
func (player *basePlayer) payout(handIdx int, result int) int {
	wager := player.Hands[handIdx].Wager
	switch result {
	case c.RESULT_BLACKJACK:
		return (wager * 3 / 2) + wager
	case c.RESULT_WIN:
		return wager + wager
	case c.RESULT_PUSH:
		return wager
	case c.RESULT_BUST, c.RESULT_LOSE:
		return 0
	}
	return 0
}

// result payout is called at the end of a turn (does not call payout if bust or blackjack)
func (player *basePlayer) resultPayout(handIdx int, result int) {
	wager := player.Hands[handIdx].Wager
	payout := player.payout(handIdx, result)
	switch result {
	case c.RESULT_WIN:
		c.Print("%s beat dealer and wins %d chips!\n", player.Name, wager)
		player.Chips += payout
	case c.RESULT_PUSH:
		c.Print("%s ties dealer and pushes.\n", player.Name)
		player.Chips += payout
	case c.RESULT_LOSE:
		// do not add payout for lose, money has already been taken
		c.Print("%s lost to dealer with a %s and loses %d chips.\n",
			player.Name,
			player.Hands[handIdx].StringSumReadable(),
			wager,
		)
	case c.RESULT_BLACKJACK:
		// do not add payout for blackjack, money has already been given
		c.Print("%s had a blackjack and earned %d chips.\n", player.Name, payout-wager)
	case c.RESULT_BUST:
		// do not add payout for bust, money has already been taken
		c.Print("%s busted and lost %d chips.\n", player.Name, wager)
	}
}

// Reset -----------------------------------------------------------------------------------

func (player *basePlayer) Reset(minBet int) {
	player.Hands = []*cards.Hand{cards.NewHand()}
	player.handsPlayed++
	if player.Chips >= minBet {
		player.Status = c.PLAYER_READY
	} else {
		player.LeaveSeat()
	}
}

// Leave -----------------------------------------------------------------------------------

func (player *basePlayer) LeaveSeat() {
	c.Print("%s has left with %d chips.\n", player.Name, player.Chips)
	player.Status = c.PLAYER_OUT
}

// Summarize -----------------------------------------------------------------------------------

func (player *basePlayer) Summarize() *stats.Stats {
	c.Print("%s has %d chips after %d hands.\n", player.Name, player.Chips, player.handsPlayed)
	return stats.NewStats(player.Name, player.Chips-player.playerRules.StartingChips, player.handsPlayed)
}

// HELPERS -----------------------------------------------------------------------------------------

// Hands
func (player *basePlayer) GetHands() []*cards.Hand {
	return player.Hands
}

func (player *basePlayer) GetHand(handIdx int) *cards.Hand {
	if len(player.Hands) <= handIdx {
		return nil
	}
	return player.Hands[handIdx]
}

func (player *basePlayer) GetChips() int {
	return player.Chips
}

// StatusIs returns true if status is one of the strings
func (player *basePlayer) StatusIs(statuses ...int) bool {
	for _, status := range statuses {
		if status == player.Status {
			return true
		}
	}
	return false
}

// PrintVisualHand prints the hand in shap of a card
func (player *basePlayer) PrintVisualHand(handIdx int) {
	c.Print("\n= %s's Hand =\n", player.Name)
	c.Print("Human has a %s.\n", player.Hands[handIdx].StringSumReadable())
	c.Print("%s\n", player.Hands[handIdx].StringLongformReadable())
}
