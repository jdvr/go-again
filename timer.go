package again

import "time"

type clock interface {
	Now() time.Time
}

type systemClock struct{}

func (sc systemClock) Now() time.Time {
	return time.Now()
}

type defaultTimer struct {
	timer *time.Timer
}

func (t *defaultTimer) Wait() <-chan time.Time {
	return t.timer.C
}

func (t *defaultTimer) Start(tick Tick) {
	if t.timer == nil {
		t.timer = time.NewTimer(tick.Next)
	} else {
		t.timer.Reset(tick.Next)
	}
}

// Stop is called when the timer is not used anymore and resources may be freed.
func (t *defaultTimer) Stop() {
	if t.timer != nil {
		t.timer.Stop()
	}
}
