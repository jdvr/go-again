package internal

import (
	"math/rand"
	"time"
)

// go port for https://github.com/googleapis/google-http-java-client/blob/da1aa993e90285ec18579f1553339b00e19b3ab5/google-http-client/src/main/java/com/google/api/client/util/ExponentialBackOff.java#L20
// Refer to original documentation for more details

const (
	// The random factor allows the code to generate values +/- 50% of the expected interval
	randomizationFactor = 0.5

	defaultInitialInterval = 500 * time.Millisecond
	defaultMaxInterval     = 30 * time.Second
	defaultMultiplier      = 1.5
	defaultTimeout         = 1 * time.Minute
)

//clock is a time wrapper
type clock interface {
	Now() time.Time
}

// BackoffConfiguration Set values for backoff algorithm configurable parameters.
type BackoffConfiguration struct {
	// InitialInterval delay before the first retry
	InitialInterval time.Duration
	// MaxInterval delay between retries, once it reaches it stop increasing
	MaxInterval time.Duration
	// IntervalMultiplier set the base interval to multiply the period delay
	IntervalMultiplier float64
	// Timeout define the max duration of the retry process
	Timeout time.Duration
	// DisableRandomization generate predicable exponential backoff intervals
	DisableRandomization bool
}

type exponentialBackoffTicksCalculator struct {
	Configuration BackoffConfiguration

	currentDelay time.Duration
	startTime    time.Time

	clock clock
}

var _ TicksCalculator = &exponentialBackoffTicksCalculator{}

func MustExponentialBackoffTicksCalculator(configuration BackoffConfiguration, clock clock) *exponentialBackoffTicksCalculator {
	return &exponentialBackoffTicksCalculator{
		Configuration: fillWithDefault(configuration),
		startTime:     clock.Now(),
		clock:         clock,
	}

}

func fillWithDefault(configuration BackoffConfiguration) BackoffConfiguration {
	initialInterval := configuration.InitialInterval
	if initialInterval == 0 {
		initialInterval = defaultInitialInterval
	}
	maxInterval := configuration.MaxInterval
	if maxInterval == 0 {
		maxInterval = defaultMaxInterval
	}
	intervalMultiplier := configuration.IntervalMultiplier
	if intervalMultiplier == 0 {
		intervalMultiplier = defaultMultiplier
	}
	timeout := configuration.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	return BackoffConfiguration{
		InitialInterval:      initialInterval,
		MaxInterval:          maxInterval,
		IntervalMultiplier:   intervalMultiplier,
		Timeout:              timeout,
		DisableRandomization: configuration.DisableRandomization,
	}
}

// Next calculates the next delay interval for a retry using currentDelay
// if DisableRandomization is true it just use IntervalMultiplier
// otherwise it use randomizationFactor to generate a "random" delta and chose a value between the min and the max
// [lastDelay - randomDelta, lastDelay + randomDelta]
// random delta is the result of randomizationFactor * currentDelay
func (c *exponentialBackoffTicksCalculator) Next() Tick {
	elapsed := c.clock.Now().Sub(c.startTime)

	var next time.Duration
	if c.Configuration.DisableRandomization {
		next = c.nextDelay()
	} else {
		current := c.currentDelay
		if current == 0 {
			current = c.Configuration.InitialInterval
		}
		next = getRandomValueFromInterval(randomizationFactor, rand.Float64(), current)
	}

	c.currentDelay = c.nextDelay()

	if elapsed > c.Configuration.Timeout {
		return Tick{
			Stop: true,
		}
	}

	return Tick{
		Next: next,
		Stop: false,
	}
}

// nextDelay generate a delay of current delay using multiplier without overflow.
func (c *exponentialBackoffTicksCalculator) nextDelay() time.Duration {
	if c.currentDelay == 0 {
		return c.Configuration.InitialInterval
	}

	next := time.Duration(float64(c.currentDelay) * c.Configuration.IntervalMultiplier)
	if next > c.Configuration.MaxInterval {
		return c.Configuration.MaxInterval
	}

	return next
}

func (c *exponentialBackoffTicksCalculator) Reset() {
	c.startTime = c.clock.Now()
	c.currentDelay = 0
}

// getRandomValueFromInterval returns a random value from the interval [randomizationFactor * currentInterval,
// * randomizationFactor * currentInterval].
func getRandomValueFromInterval(randomizationFactor, random float64, currentInterval time.Duration) time.Duration {
	var delta = randomizationFactor * float64(currentInterval)
	var minInterval = float64(currentInterval) - delta
	var maxInterval = float64(currentInterval) + delta

	// Get a random value from the range [minInterval, maxInterval].
	// The formula used below has a +1 because if the minInterval is 1 and the maxInterval is 3 then
	// we want a 33% chance for selecting either 1, 2 or 3.
	return time.Duration(minInterval + (random * (maxInterval - minInterval + 1)))
}
