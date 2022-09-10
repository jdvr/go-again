package again

import (
	"context"
	"time"
)

func WithExponentialBackoff(configuration BackoffConfiguration) Retryer {
	return mustDefaultRetryerWithDefaultTimer(retryerConfig{
		TicksCalculator: newExponentialBackoffTicksCalculator(configuration),
	})
}

func WithConstantDelay(delay, timeout time.Duration) Retryer {
	return mustDefaultRetryerWithDefaultTimer(retryerConfig{
		TicksCalculator: mustConstantDelayTicksCalculator(delay, timeout),
	})
}

func WithCustomTicksCalculator(calculator TicksCalculator) Retryer {
	return mustDefaultRetryerWithDefaultTimer(retryerConfig{
		TicksCalculator: calculator,
	})
}

func Retry(ctx context.Context, operation Operation) error {
	retryer := WithExponentialBackoff(BackoffConfiguration{})

	err := retryer.Retry(ctx, operation)
	if err != nil {
		return err
	}

	return nil
}
