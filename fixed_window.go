package rate_limiter

import (
	"sync"
	"time"
)

type FWLimiter struct {
	requestsLimit int
	window        time.Duration
	requests      int
	reset         time.Time
	mu            sync.Mutex
}

func NewFWLimiter(requestsLimit int, window time.Duration) *FWLimiter {
	return &FWLimiter{
		requestsLimit: requestsLimit,
		window:        window,
	}
}

func (fwl *FWLimiter) Allow() bool {
	fwl.mu.Lock()
	defer fwl.mu.Unlock()

	current := time.Now()

	if fwl.reset.IsZero() {
		fwl.reset = current.Add(fwl.window)
	}

	if current.After(fwl.reset) {
		fwl.requests = 0
		fwl.reset = current.Add(fwl.window)
	}

	if fwl.requests < fwl.requestsLimit {
		fwl.requests++
		return true
	}

	return false
}
