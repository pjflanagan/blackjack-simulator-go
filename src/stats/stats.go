package stats

// TODO: make this a general class so all players and dealer can record stats
// learners can record scenario stats
// all players can record house odds on a typical hand
// we can start doing dealer scenario's too eventually
// https://math.info/Misc/House_Edge/
// potentailGain * oddsOfWin + potentialLoss * oddsOfLoss

// if we make the players all play a max number of hands then we could make
// the odds more accurate, otherwise we'd just be averaging 100%'s and -100%'s

type Stats struct {
	strategy string // type of
	dChips   int    // delta chips
}

func NewStats(strategy string, dChips int) *Stats {
	return &Stats{
		strategy: strategy,
		dChips:   dChips,
	}
}

func (stats *Stats) Exchange(dChips int) {
	// we need this function so the learner, who never wins or loses any chips
	// can record the amount of chips maybe idk
	stats.dChips += dChips
}

// in order for stats to work we need to play many games, no one game
// can be used to compute house edge, we need to set a max hands
// for each other players.
// if we sum all of the positive gain amounts and multiply that by the times it happens
