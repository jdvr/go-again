package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExponentialBackoffTicksCalculator_Next(t *testing.T) {
	t.Run("linear backoff respecting the max interval", func(t *testing.T) {
		ticksCalculator := MustExponentialBackoffTicksCalculator(BackoffConfiguration{
			InitialInterval:      500 * time.Millisecond,
			MaxInterval:          5 * time.Second,
			IntervalMultiplier:   2,
			Timeout:              10 * time.Second,
			DisableRandomization: true,
		}, defaultClock{})
		expected := []Tick{
			{Next: 500 * time.Millisecond},
			{Next: 1000 * time.Millisecond},
			{Next: 2000 * time.Millisecond},
			{Next: 4000 * time.Millisecond},
			{Next: 5000 * time.Millisecond},
			{Next: 5000 * time.Millisecond},
		}

		var generated []Tick
		for i := 0; i < len(expected); i++ {
			generated = append(generated, ticksCalculator.Next())
		}

		require.Equal(t, expected, generated)
	})
	t.Run("stop when timed out", func(t *testing.T) {
		ticksCalculator := MustExponentialBackoffTicksCalculator(BackoffConfiguration{
			Timeout:              1 * time.Nanosecond,
			DisableRandomization: true,
		}, defaultClock{})

		require.Equal(t, Tick{Stop: true}, ticksCalculator.Next())
	})
	t.Run("generate random values for intervals", func(t *testing.T) {
		ticksCalculator := MustExponentialBackoffTicksCalculator(BackoffConfiguration{
			InitialInterval:    500 * time.Millisecond,
			MaxInterval:        5 * time.Second,
			IntervalMultiplier: 2,
			Timeout:            10 * time.Second,
		}, defaultClock{})
		expected := []Tick{
			{Next: 500 * time.Millisecond},
			{Next: 1000 * time.Millisecond},
			{Next: 2000 * time.Millisecond},
			{Next: 4000 * time.Millisecond},
			{Next: 5000 * time.Millisecond},
			{Next: 5000 * time.Millisecond},
		}

		var generated []Tick
		for i := 0; i < len(expected); i++ {
			generated = append(generated, ticksCalculator.Next())
		}

		assertProgressiveValues(t, generated)

	})
	t.Run("configuration is fill with default values", func(t *testing.T) {
		defaultConfiguration := MustExponentialBackoffTicksCalculator(BackoffConfiguration{}, defaultClock{}).Configuration

		require.Equal(t, BackoffConfiguration{
			InitialInterval:    defaultInitialInterval,
			MaxInterval:        defaultMaxInterval,
			IntervalMultiplier: defaultMultiplier,
			Timeout:            defaultTimeout,
		}, defaultConfiguration)
	})
}

func assertProgressiveValues(t *testing.T, generated []Tick) {
	for i := 1; i < len(generated); i++ {
		require.Less(t, generated[i-1].Next, generated[i].Next)
	}
}
