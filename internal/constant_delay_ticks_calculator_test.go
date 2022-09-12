package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConstantDelayTicksCalculator_Next(t *testing.T) {
	t.Run("delay is always the same", func(t *testing.T) {
		ticksCalculator := MustConstantDelayTicksCalculator(500*time.Millisecond, 1*time.Hour, defaultClock{})
		expected := []Tick{
			{Next: 500 * time.Millisecond},
			{Next: 500 * time.Millisecond},
			{Next: 500 * time.Millisecond},
			{Next: 500 * time.Millisecond},
		}

		var generated []Tick
		for i := 0; i < len(expected); i++ {
			generated = append(generated, ticksCalculator.Next())
		}

		require.Equal(t, expected, generated)
	})
	t.Run("stop when timed out", func(t *testing.T) {
		ticksCalculator := MustConstantDelayTicksCalculator(500*time.Millisecond, 1*time.Nanosecond, defaultClock{})

		require.Equal(t, Tick{Stop: true}, ticksCalculator.Next())
	})
	t.Run("panics for 0 config", func(t *testing.T) {
		require.Panics(t, func() {
			MustConstantDelayTicksCalculator(0, time.Second, defaultClock{})
		})
		require.Panics(t, func() {
			MustConstantDelayTicksCalculator(time.Hour, 0, defaultClock{})
		})
	})
}
