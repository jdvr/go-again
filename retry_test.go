package again_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/jdvr/go-again"
)

func TestRetryer_Retry(t *testing.T) {
	t.Parallel()

	t.Run("stop if next tick includes stop flag", func(t *testing.T) {
		t.Parallel()

		givenFakeOperation := NewFakeOperation(t)
		givenCtx := context.TODO()

		givenFakeOperation.
			givenContext(givenCtx).
			Returns(nil)

		retrayer := again.mustRetryer(again.retryerConfig{
			TicksCalculator: singleTicksCalculator{},
			Timer:           &instantTimer{},
		})

		err := retrayer.Retry(givenCtx, givenFakeOperation)

		require.NoError(t, err)
		givenFakeOperation.haveBeenCalled(1)
	})

	t.Run("stop if operation if error is permanent and return underlying error", func(t *testing.T) {
		t.Parallel()

		givenFakeOperation := NewFakeOperation(t)
		givenCtx := context.TODO()

		expectedError := errors.New("whatever")

		givenFakeOperation.
			givenContext(givenCtx).
			Returns(again.Permanent(errors.New("whatever")))

		retrayer := again.mustRetryer(again.retryerConfig{
			TicksCalculator: singleTicksCalculator{},
			Timer:           &instantTimer{},
		})

		err := retrayer.Retry(givenCtx, givenFakeOperation)

		require.Equal(t, expectedError, err)
		givenFakeOperation.haveBeenCalled(1)
	})
	t.Run("stop whenever context is cancelled", func(t *testing.T) {
		t.Parallel()

		givenFakeOperation := NewFakeOperation(t)
		givenCtx, cancel := context.WithCancel(context.TODO())

		givenFakeOperation.
			givenContext(givenCtx).
			Returns(nil)

		retrayer := again.mustRetryer(again.retryerConfig{
			TicksCalculator: infinityTicksCalculator{},
			Timer:           &instantTimer{},
		})

		cancel()
		err := retrayer.Retry(givenCtx, givenFakeOperation)

		require.Equal(t, context.Canceled, err)
		givenFakeOperation.haveBeenCalled(1)
	})
	t.Run("operation is executed while ticks are returned", func(t *testing.T) {
		t.Parallel()

		givenFakeOperation := NewFakeOperation(t)
		givenCtx := context.TODO()

		givenFakeOperation.
			givenContext(givenCtx).
			Returns(nil)

		retrayer := again.mustRetryer(again.retryerConfig{
			TicksCalculator: &twoTicksCalculator{},
			Timer:           &instantTimer{},
		})

		err := retrayer.Retry(givenCtx, givenFakeOperation)

		require.NoError(t, err)
		givenFakeOperation.haveBeenCalled(2)
	})
	t.Run("operation is executed while ticks are returned event with error", func(t *testing.T) {
		t.Parallel()

		givenFakeOperation := NewFakeOperation(t)
		givenCtx := context.TODO()

		anyError := errors.New("any error")

		givenFakeOperation.
			givenContext(givenCtx).
			Returns(anyError)

		retrayer := again.mustRetryer(again.retryerConfig{
			TicksCalculator: &twoTicksCalculator{},
			Timer:           &instantTimer{},
		})

		err := retrayer.Retry(givenCtx, givenFakeOperation)

		require.ErrorIs(t, err, anyError)
		givenFakeOperation.haveBeenCalled(2)
	})
}

func TestPermanentError(t *testing.T) {
	t.Parallel()

	t.Run("unwrap", func(t *testing.T) {
		t.Parallel()

		givenError := errors.New("foo")
		var permanentError error = again.Permanent(givenError)

		unWrapped := errors.Unwrap(permanentError)

		require.Equal(t, givenError, unWrapped)
	})

	t.Run("Is", func(t *testing.T) {
		t.Parallel()

		givenError := errors.New("given")
		givenOtherError := errors.New("other")

		var permanentError error = again.Permanent(givenError)

		require.ErrorIs(t, permanentError, givenError)
		require.NotErrorIs(t, permanentError, givenOtherError)
	})

	t.Run("As", func(t *testing.T) {
		t.Parallel()

		givenError := errors.New("given")

		var permanentError error = again.Permanent(givenError)

		var permanentRef *again.PermanentError

		require.ErrorAs(t, permanentError, &permanentRef)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		givenError := errors.New("given error")

		var permanentError error = again.Permanent(givenError)

		require.Equal(t, "given error", permanentError.Error())
	})
}

type instantTimer struct {
	timer *time.Timer
}

func (i *instantTimer) Start(_ again.Tick) {
	i.timer = time.NewTimer(1 * time.Nanosecond)
}

func (i *instantTimer) Wait() <-chan time.Time {
	return i.timer.C
}

func (i *instantTimer) Stop() {
	// Someone could try to stop a timer before start it.
	if i.timer != nil {
		i.timer.Stop()
	}
}

type singleTicksCalculator struct{}

func (s singleTicksCalculator) Next() again.Tick {
	return again.Tick{Next: -1, Stop: true}
}

func (s singleTicksCalculator) Reset() {}

type infinityTicksCalculator struct{}

func (s infinityTicksCalculator) Next() again.Tick {
	return again.Tick{Next: 100 * time.Hour, Stop: false}
}

func (s infinityTicksCalculator) Reset() {}

type twoTicksCalculator struct {
	called int
}

func (ticksCalculator *twoTicksCalculator) Next() again.Tick {
	ticksCalculator.called += 1
	if ticksCalculator.called == 2 {
		return again.Tick{Next: 1 * time.Hour, Stop: true}
	}
	return again.Tick{Next: 1 * time.Millisecond, Stop: false}
}

func (ticksCalculator *twoTicksCalculator) Reset() {
	ticksCalculator.called = 0
}
