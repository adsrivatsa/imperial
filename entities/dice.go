package entities

import (
	"math"
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
	// and uses exponential weighting so that more-rolled faces become less
	// likely, pushing the distribution toward uniformity over time.
	WeightedDie struct {
		// Counts[i] is the number of times face i+1 has been rolled.
		Counts [6]int `msgpack:"c"`
		// Alpha controls the strength of the bias.
		// Higher values make the correction more aggressive.
		Alpha float64 `msgpack:"a"`
	}
)

// NewWeightedDie creates a WeightedDie with the given alpha (decay rate).
// A reasonable default alpha is 0.3 – strong enough to smooth the
// distribution without making it completely deterministic.
func NewWeightedDie(alpha float64) *WeightedDie {
	return &WeightedDie{
		Counts: [6]int{},
		Alpha:  alpha,
	}
}

// Roll returns a value in [1, 6] chosen with exponential weighting.
// Each face i has weight  w_i = exp(-alpha * count_i).
// The probability of rolling face i is  w_i / sum(w_j).
func (d *WeightedDie) Roll() int {
	var weights [6]float64
	var total float64

	for i := 0; i < 6; i++ {
		weights[i] = math.Exp(-d.Alpha * float64(d.Counts[i]))
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
