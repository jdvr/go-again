package internal

import (
	"time"
)

type constantDelayTicksCalculator struct {
	delay   time.Duration
	timeout time.Duration

	startAt time.Time
	clock   clock
}

func MustConstantDelayTicksCalculator(delay time.Duration, timeout time.Duration, clock clock) TicksCalculator {
	if delay == 0 || timeout == 0 {
		panic("delay and timeout must be set")
	}
	return &constantDelayTicksCalculator{
		delay:   delay,
		timeout: timeout,
		startAt: clock.Now(),
		clock:   clock,
	}
}

func (c *constantDelayTicksCalculator) Next() Tick {
	elapsed := c.clock.Now().Sub(c.startAt)
	if elapsed > c.timeout {
		return Tick{Stop: true}
	}

	return Tick{
		Next: c.delay,
		Stop: false,
	}
}

func (c *constantDelayTicksCalculator) Reset() {
	c.startAt = c.clock.Now()
}
