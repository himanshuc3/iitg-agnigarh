// NOTE: Shamelessly copied from
// https://chromium.googlesource.com/external/github.com/cenkalti/backoff/+/refs/tags/v2.1.0/backoff.go
// NOTE: Need to create a Ticker with function context to
// execute it on a loop

package utils

import (
	"math/rand"
	"time"
)

const Stop time.Duration = -1

/*

	NextBackoff = currentInterval + random * (max-min)

*/

// NOTE: Technically some of the fields could be hardcoded
type ExponentialBackoff struct {
	// For retry
	InitialInterval time.Duration
	MaxInterval     time.Duration
	// Stop at some point (0 retrying beyond infinity)
	MaxElapsedTime      time.Duration
	RandomizationFactor float64
	// Next interval would get this time
	Multiplier float64

	// Not required on instance, can be global
	// Clock Clock

	currentInterval time.Duration
	startTime       time.Time
}

// Clock is an interface that returns current time for BackOff
// Not used to using interfaces
// Any type implementing methods from interface
// can be disguised as that interace
// type Clock interface {
// 	Now() time.Time
// }

// Default values for ExponentialBackoff
const (
	DefaultInitialInterval     = 2000 * time.Millisecond
	DefaultRandomizationFactor = 0.5
	DefaultMultiplier          = 1.5
	DefaultMaxInterval         = 60 * time.Second
	DefaultMaxElapsedTime      = 10 * time.Minute
)

func NewExponentialBackoff() *ExponentialBackoff {
	b := &ExponentialBackoff{
		InitialInterval:     DefaultInitialInterval,
		RandomizationFactor: DefaultRandomizationFactor,
		Multiplier:          DefaultMultiplier,
		MaxInterval:         DefaultMaxInterval,
		MaxElapsedTime:      DefaultMaxElapsedTime,
	}
	b.Reset()
	return b

}

// NOTE: naming this would still be better than one variable names tbh
func (b *ExponentialBackoff) Reset() {
	b.currentInterval = b.InitialInterval
	b.startTime = time.Now()
}

func (b *ExponentialBackoff) NextBackoff() time.Duration {
	if b.MaxElapsedTime != 0 && b.GetElapsedTime() > b.MaxElapsedTime {
		return Stop
	}

	defer b.incrementCurrentInterval()

	return getRandomValueFromInterval(b.RandomizationFactor, b.currentInterval)
}

func (b *ExponentialBackoff) GetElapsedTime() time.Duration {
	return time.Since(b.startTime)
}

func (b *ExponentialBackoff) incrementCurrentInterval() {
	if float64(b.currentInterval)*b.Multiplier >= float64(b.MaxInterval) {
		b.currentInterval = b.MaxInterval
	} else {
		b.currentInterval = time.Duration(float64(b.currentInterval) * b.Multiplier)
	}

}

func getRandomValueFromInterval(randomizationFactor float64, currentInterval time.Duration) time.Duration {
	var delta = randomizationFactor * float64(currentInterval)
	var minInterval = float64(currentInterval) - delta
	var maxInterval = float64(currentInterval) + delta

	// The formula used below has a +1 because if the minInterval is 1 and the maxInterval is 3 then
	// we want a 33% chance for selecting either 1, 2 or 3.
	return time.Duration(minInterval + (rand.Float64() * (maxInterval - minInterval + 1)))
}
