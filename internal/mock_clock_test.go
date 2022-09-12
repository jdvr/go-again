package internal

import (
	"time"
)

type defaultClock struct{}

func (d defaultClock) Now() time.Time {
	return time.Now()
}
