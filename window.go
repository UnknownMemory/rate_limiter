package rate_limiter

import (
	"time"
)

type Window struct {
	start time.Time
	count int
}

func NewWindow(start time.Time) *Window {
	return &Window{
		start: start,
		count: 0,
	}
}

func (w *Window) Set(start time.Time, count int) {
	w.start = start
	w.count = count
}
