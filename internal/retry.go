package internal

import (
	"context"
	"errors"
	"time"
)

type Operation interface {
	Run(context context.Context) error
}

type Tick struct {
	Next time.Duration
	Stop bool
}

type Timer interface {
	Start(tick Tick)
	Wait() <-chan time.Time
	Stop()
}

type TicksCalculator interface {
	Next() Tick
	Reset()
}

type Retryer interface {
	Retry(ctx context.Context, operation Operation) error
}

type defaultRetryer struct {
	TicksCalculator TicksCalculator
	Timer           Timer
}

type RetryerConfig struct {
	TicksCalculator TicksCalculator
	Timer           Timer
}

// MustRetryer returns a new Retryer or panic if any dependency is nil.
func MustRetryer(config RetryerConfig) Retryer {
	if config.Timer == nil {
		panic("again: MustRetryer: nil Timer")
	}
	if config.TicksCalculator == nil {
		panic("again: MustRetryer: nil TicksCalculator")
	}
	return defaultRetryer{
		TicksCalculator: config.TicksCalculator,
		Timer:           config.Timer,
	}
}

func (retryer defaultRetryer) Retry(ctx context.Context, operation Operation) error {
	var next Tick

	defer func() {
		retryer.Timer.Stop()
	}()

	retryer.TicksCalculator.Reset()
	for {
		err := operation.Run(ctx)

		var permanent *PermanentError
		if errors.As(err, &permanent) {
			return permanent.Err
		}

		if next = retryer.TicksCalculator.Next(); next.Stop {
			if cerr := ctx.Err(); cerr != nil {
				return cerr
			}
			return err
		}

		retryer.Timer.Start(next)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-retryer.Timer.Wait():
		}
	}
}

type PermanentError struct {
	Err error
}

func (e *PermanentError) Error() string {
	return e.Err.Error()
}

func (e *PermanentError) Unwrap() error {
	return e.Err
}

func (e *PermanentError) Is(target error) bool {
	_, ok := target.(*PermanentError)
	return ok
}

func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return &PermanentError{
		Err: err,
	}
}
