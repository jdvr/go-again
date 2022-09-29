package again

import (
	"context"
	"time"

	"github.com/jdvr/go-again/internal"
)

type Operation[T any] internal.Operation[T]
type BackoffConfiguration = internal.BackoffConfiguration

// WithExponentialBackoff initialize a retryer using ExponentialBackoff algorithm to calculate delay between each retry.
func WithExponentialBackoff[T any](configuration BackoffConfiguration) internal.Retryer[T] {
	return internal.MustRetryer[T](internal.RetryerConfig{
		TicksCalculator: internal.MustExponentialBackoffTicksCalculator(configuration, systemClock{}),
		Timer:           &defaultTimer{},
	})
}

// WithConstantDelay initialize a retryer using a constant delay algorithm to calculate delay between each retry.
func WithConstantDelay[T any](delay, timeout time.Duration) internal.Retryer[T] {
	return internal.MustRetryer[T](internal.RetryerConfig{
		TicksCalculator: internal.MustConstantDelayTicksCalculator(delay, timeout, systemClock{}),
		Timer:           &defaultTimer{},
	})
}

// WithCustomTicksCalculator initialize retryer using custom calculator to calculate delays between retries.
func WithCustomTicksCalculator[T any](calculator internal.TicksCalculator) internal.Retryer[T] {
	return internal.MustRetryer[T](internal.RetryerConfig{
		TicksCalculator: calculator,
		Timer:           &defaultTimer{},
	})
}

// RetryOperation use WithExponentialBackoff to retry operation until it stops failing or timeout is reached.
// it might return the last operation run error or a context cancelled error.
func RetryOperation[T any](ctx context.Context, operation Operation[T]) (T, error) {
	retryer := WithExponentialBackoff[T](internal.BackoffConfiguration{})

	value, err := retryer.Retry(ctx, operation)
	if err != nil {
		return nil, err
	}

	return value, nil
}

type RunFunc[T any] func(context.Context) (T, error)

// Retry use WithExponentialBackoff to retry the run function until it stops failing or timeout is reached.
// it might return the last function run error or a context cancelled error.
func Retry[T any](ctx context.Context, run RunFunc[T]) (T, error) {
	return RetryOperation[T](ctx, handleRun(run))
}

type wrappedRun[T any] struct {
	run RunFunc[T]
}

func (w wrappedRun[T]) Run(context context.Context) (T, error) {
	return w.run(context)
}

func handleRun[T any](run RunFunc[T]) internal.Operation[T] {
	return wrappedRun[T]{run: run}
}

func NewPermanentError(err error) error {
	if err == nil {
		return nil
	}
	return &internal.PermanentError{
		Err: err,
	}
}
