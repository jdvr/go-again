package again

import (
	"context"
	"github.com/jdvr/go-again/internal"
	"time"
)

type Operation = internal.Operation
type BackoffConfiguration = internal.BackoffConfiguration

func WithExponentialBackoff(configuration BackoffConfiguration) internal.Retryer {
	return internal.MustRetryer(internal.RetryerConfig{
		TicksCalculator: internal.MustExponentialBackoffTicksCalculator(configuration, systemClock{}),
		Timer:           &defaultTimer{},
	})
}

func WithConstantDelay(delay, timeout time.Duration) internal.Retryer {
	return internal.MustRetryer(internal.RetryerConfig{
		TicksCalculator: internal.MustConstantDelayTicksCalculator(delay, timeout, systemClock{}),
		Timer:           &defaultTimer{},
	})
}

func WithCustomTicksCalculator(calculator internal.TicksCalculator) internal.Retryer {
	return internal.MustRetryer(internal.RetryerConfig{
		TicksCalculator: calculator,
		Timer:           &defaultTimer{},
	})
}

func RetryOperation(ctx context.Context, operation Operation) error {
	retryer := WithExponentialBackoff(internal.BackoffConfiguration{})

	err := retryer.Retry(ctx, operation)
	if err != nil {
		return err
	}

	return nil
}

type RunFunc func(context.Context) error

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
