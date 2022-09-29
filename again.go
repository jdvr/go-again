package again

import (
	"context"
	"time"

	"github.com/jdvr/go-again/internal"
)

type Operation = internal.Operation
type BackoffConfiguration = internal.BackoffConfiguration

// WithExponentialBackoff initialize a retryer using ExponentialBackoff algorithm to calculate delay between each retry.
func WithExponentialBackoff(configuration BackoffConfiguration) internal.Retryer {
	return internal.MustRetryer(internal.RetryerConfig{
		TicksCalculator: internal.MustExponentialBackoffTicksCalculator(configuration, systemClock{}),
		Timer:           &defaultTimer{},
	})
}

// WithConstantDelay initialize a retryer using a constant delay algorithm to calculate delay between each retry.
func WithConstantDelay(delay, timeout time.Duration) internal.Retryer {
	return internal.MustRetryer(internal.RetryerConfig{
		TicksCalculator: internal.MustConstantDelayTicksCalculator(delay, timeout, systemClock{}),
		Timer:           &defaultTimer{},
	})
}

// WithCustomTicksCalculator initialize retryer using custom calculator to calculate delays between retries.
func WithCustomTicksCalculator(calculator internal.TicksCalculator) internal.Retryer {
	return internal.MustRetryer(internal.RetryerConfig{
		TicksCalculator: calculator,
		Timer:           &defaultTimer{},
	})
}

// RetryOperation use WithExponentialBackoff to retry operation until it stops failing or timeout is reached.
// it might return the last operation run error or a context cancelled error.
func RetryOperation(ctx context.Context, operation Operation) error {
	retryer := WithExponentialBackoff(internal.BackoffConfiguration{})

	err := retryer.Retry(ctx, operation)
	if err != nil {
		return err
	}

	return nil
}

type RunFunc func(context.Context) error

// Retry use WithExponentialBackoff to retry the run function until it stops failing or timeout is reached.
// it might return the last function run error or a context cancelled error.
func Retry(ctx context.Context, run RunFunc) error {
	return RetryOperation(ctx, handleRun(run))
}

type wrappedRun struct {
	run RunFunc
}

func (w wrappedRun) Run(context context.Context) error {
	return w.run(context)
}

func handleRun(run RunFunc) internal.Operation {
	return wrappedRun{run: run}
}

func NewPermanentError(err error) error {
	if err == nil {
		return nil
	}
	return &internal.PermanentError{
		Err: err,
	}
}
