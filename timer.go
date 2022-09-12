package again

import (
	"github.com/jdvr/go-again/internal"
	"time"
)

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

func (t *defaultTimer) Start(tick internal.Tick) {
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
