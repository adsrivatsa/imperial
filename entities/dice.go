package entities

import (
	"math/rand"
)

type (
	DieRollState struct {
		RedRoll          int            `msgpack:"r"`
		WhiteRoll        int            `msgpack:"w"`
		EventRoll        int            `msgpack:"e"`
		PlayerHandDeltas []*Hand        `msgpack:"-"`
		GainInfo         []CardMoveInfo `msgpack:"g"`
		IsInit           bool           `msgpack:"ii,omitempty"`
	}

	DiceStats struct {
		Rolls      [12]int `msgpack:"r"`
		EventRolls [6]int  `msgpack:"e"`
	}

	// WeightedDie tracks the number of times each face (1-6) has been rolled
	// and uses inverse-probability weighting so that more-rolled faces become
	// less likely, pushing the distribution toward uniformity over time.
	WeightedDie struct {
		// Counts[i] is the number of times face i+1 has been rolled.
		Counts [6]int `msgpack:"c"`
	}
)

// NewWeightedDie creates a new WeightedDie with zeroed counts.
func NewWeightedDie() *WeightedDie {
	return &WeightedDie{
		Counts: [6]int{},
	}
}

// Roll returns a value in [1, 6] chosen with inverse-probability weighting.
// Each face i has weight  w_i = 1 / (1 + count_i).
// Faces rolled more often receive proportionally less weight.
func (d *WeightedDie) Roll() int {
	var weights [6]float64
	var total float64

	for i := 0; i < 6; i++ {
		weights[i] = 1.0 / (1.0 + float64(d.Counts[i]))
		total += weights[i]
	}

	// Pick a random value in [0, total)
	r := rand.Float64() * total
	var cumulative float64
	for i := 0; i < 6; i++ {
		cumulative += weights[i]
		if r < cumulative {
			d.Counts[i]++
			return i + 1
		}
	}

	// Fallback (should not normally be reached due to floating point)
	d.Counts[5]++
	return 6
}

type GameOverMessage struct {
	Players []*PlayerState `msgpack:"p"`
	Winner  uint16         `msgpack:"w"`
}
