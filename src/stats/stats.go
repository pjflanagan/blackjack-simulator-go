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
	strategy    string // type of
	dChips      int    // change in chips over hands played
	handsPlayed int    // number of hands played
}

func NewStats(strategy string, dChips int, handsPlayed int) *Stats {
	return &Stats{
		strategy:    strategy,
		dChips:      dChips,
		handsPlayed: handsPlayed,
	}
}

func (stats *Stats) GetStrategy() string {
	return stats.strategy
}

func (stats *Stats) ExpectedGain() float32 {
	return float32(stats.dChips) / float32(stats.handsPlayed)
}

// in order for stats to work we need to play many games, no one game
// can be used to compute house edge, we need to set a max hands
// for each other players.
// if we sum all of the positive gain amounts and multiply that by the times it happens
func HouseOdds(allStats []*Stats) map[string]*Stats {
	finalStats := make(map[string]*Stats)
	for _, stat := range allStats {
		if finalStats[stat.strategy] == nil {
			finalStats[stat.strategy] = stat
		} else {
			finalStats[stat.strategy].dChips += stat.dChips
			finalStats[stat.strategy].handsPlayed += stat.handsPlayed
		}
	}
	return finalStats
}
