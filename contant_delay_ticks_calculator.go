package again

import (
	"time"
)

type constantDelayTicksCalculator struct {
	delay time.Duration
}

func newConstantDelayTicksCalculator(delay time.Duration) TicksCalculator {
	return constantDelayTicksCalculator{
		delay: 0,
	}
}

func (c constantDelayTicksCalculator) Next() Tick {
	return Tick{
		Next: c.delay,
		Stop: false,
	}
}

func (c constantDelayTicksCalculator) Reset() {}
