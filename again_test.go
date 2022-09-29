package again_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/jdvr/go-again"
)

func TestWithExponentialBackoff(t *testing.T) {
	t.Run("given operation is called until timeout", func(t *testing.T) {
		givenOperation := NewFakeOperation(t)

		givenOperation.allowAnyCall = true

		retryer := again.WithExponentialBackoff[int](again.BackoffConfiguration{
			InitialInterval:    10 * time.Millisecond,
			MaxInterval:        50 * time.Millisecond,
			IntervalMultiplier: 2,
			Timeout:            20 * time.Millisecond,
		})

		_, err := retryer.Retry(context.Background(), givenOperation)
		require.Error(t, err)

		givenOperation.haveBeenCalled(3)
	})
	t.Run("disable random delays", func(t *testing.T) {
		givenOperation := NewFakeOperation(t)

		givenOperation.allowAnyCall = true

		retryer := again.WithExponentialBackoff[int](again.BackoffConfiguration{
			InitialInterval:      10 * time.Millisecond,
			MaxInterval:          50 * time.Millisecond,
			IntervalMultiplier:   2,
			Timeout:              20 * time.Millisecond,
			DisableRandomization: true,
		})

		_, err := retryer.Retry(context.Background(), givenOperation)
		require.Error(t, err)

		givenOperation.haveBeenCalled(3)
	})
}

func TestWithConstantDelay(t *testing.T) {
	t.Run("given operation is called until timeout", func(t *testing.T) {
		givenOperation := NewFakeOperation(t)

		givenOperation.allowAnyCall = true

		retryer := again.WithConstantDelay[int](1*time.Millisecond, 3*time.Millisecond)

		_, err := retryer.Retry(context.Background(), givenOperation)
		require.Error(t, err)

		givenOperation.haveBeenCalled(4)
	})
}

func TestRetryOperation(t *testing.T) {
	t.Run("given operation is called until permanent error", func(t *testing.T) {
		testContext := context.Background()
		givenOperation := NewFakeOperation(t)

		givenOperation.
			givenContext(testContext).
			Returns(0, errors.New("no permanent"))

		givenOperation.
			givenContext(testContext).
			Returns(0, errors.New("no permanent"))

		expectedErr := errors.New("whatever")
		givenOperation.
			givenContext(testContext).
			Returns(0, again.NewPermanentError(expectedErr))

		_, err := again.RetryOperation[int](testContext, givenOperation)
		require.ErrorIs(t, err, expectedErr)

		givenOperation.haveBeenCalled(3)
	})
}

func TestRetry(t *testing.T) {
	t.Run("given function is called until permanent error", func(t *testing.T) {
		testContext := context.Background()
		expectedErr := errors.New("whatever")
		called := 0
		givenFunc := func(ctx context.Context) (bool, error) {
			called++
			if called > 3 {
				return false, again.NewPermanentError(expectedErr)
			}

			return false, errors.New("whatever")
		}

		_, err := again.Retry[bool](testContext, givenFunc)
		require.ErrorIs(t, err, expectedErr)
		require.Equal(t, 4, called)
	})
}
