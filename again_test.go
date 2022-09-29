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

		retryer := again.WithExponentialBackoff(again.BackoffConfiguration{
			InitialInterval:    10 * time.Millisecond,
			MaxInterval:        50 * time.Millisecond,
			IntervalMultiplier: 2,
			Timeout:            20 * time.Millisecond,
		})

		err := retryer.Retry(context.Background(), givenOperation)
		require.NoError(t, err)

		givenOperation.haveBeenCalled(3)
	})
	t.Run("disable random delays", func(t *testing.T) {
		givenOperation := NewFakeOperation(t)

		givenOperation.allowAnyCall = true

		retryer := again.WithExponentialBackoff(again.BackoffConfiguration{
			InitialInterval:      10 * time.Millisecond,
			MaxInterval:          50 * time.Millisecond,
			IntervalMultiplier:   2,
			Timeout:              20 * time.Millisecond,
			DisableRandomization: true,
		})

		err := retryer.Retry(context.Background(), givenOperation)
		require.NoError(t, err)

		givenOperation.haveBeenCalled(3)
	})
}

func TestWithConstantDelay(t *testing.T) {
	t.Run("given operation is called until timeout", func(t *testing.T) {
		givenOperation := NewFakeOperation(t)

		givenOperation.allowAnyCall = true

		retryer := again.WithConstantDelay(1*time.Millisecond, 3*time.Millisecond)

		err := retryer.Retry(context.Background(), givenOperation)
		require.NoError(t, err)

		givenOperation.haveBeenCalled(4)
	})
}

func TestRetryOperation(t *testing.T) {
	t.Run("given operation is called until permanent error", func(t *testing.T) {
		testContext := context.Background()
		givenOperation := NewFakeOperation(t)

		givenOperation.
			givenContext(testContext).
			Returns(nil)

		givenOperation.
			givenContext(testContext).
			Returns(nil)

		expectedErr := errors.New("whatever")
		givenOperation.
			givenContext(testContext).
			Returns(again.NewPermanentError(expectedErr))

		err := again.RetryOperation(testContext, givenOperation)
		require.ErrorIs(t, err, expectedErr)

		givenOperation.haveBeenCalled(3)
	})
}

func TestRetry(t *testing.T) {
	t.Run("given function is called until permanent error", func(t *testing.T) {
		testContext := context.Background()
		expectedErr := errors.New("whatever")
		called := 0
		givenFunc := func(ctx context.Context) error {
			called++
			if called > 3 {
				return again.NewPermanentError(expectedErr)
			}

			return nil
		}

		err := again.Retry(testContext, givenFunc)
		require.ErrorIs(t, err, expectedErr)
		require.Equal(t, 4, called)
	})
}
