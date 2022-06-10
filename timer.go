package again

import "time"

type defaultTimer struct {
	timer *time.Timer
}

func (t *defaultTimer) C() <-chan time.Time {
	return t.timer.C
}

func (t *defaultTimer) Start(duration time.Duration) {
	if t.timer == nil {
		t.timer = time.NewTimer(duration)
	} else {
		t.timer.Reset(duration)
	}
}

// Stop is called when the timer is not used anymore and resources may be freed.
func (t *defaultTimer) Stop() {
	if t.timer != nil {
		t.timer.Stop()
	}
}
