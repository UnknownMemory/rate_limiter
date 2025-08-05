package rate_limiter

import (
	"sync"
	"time"
)

type SWLimiter struct {
	requestsLimit int
	windowSize    time.Duration
	currWindow    *Window
	prevWindow    *Window
	mu            sync.Mutex
}

func NewSWLimiter(requestsLimit int, windowSize time.Duration) *SWLimiter {
	return &SWLimiter{
		requestsLimit: requestsLimit,
		windowSize:    windowSize,
		currWindow:    NewWindow(time.Now().Truncate(windowSize)),
		prevWindow:    NewWindow(time.Now().Truncate(windowSize).Add(-windowSize)),
	}
}

func (swl *SWLimiter) Allow() bool {
	swl.mu.Lock()
	defer swl.mu.Unlock()

	now := time.Now()
	swl.updateWindows(now)

	timeElapsed := float64(now.Sub(swl.currWindow.start))
	size := float64(swl.windowSize)
	weight := (size - timeElapsed) / size
	limitApproximation := int(float64(swl.prevWindow.count)*weight) + swl.currWindow.count

	if limitApproximation < swl.requestsLimit {
		swl.currWindow.count += 1
		return true
	}

	return false
}

func (swl *SWLimiter) updateWindows(now time.Time) {
	current := now.Truncate(swl.windowSize)
	currentWinStart := swl.currWindow.start

	if !current.Equal(currentWinStart) {
		prevWinCount := 0

		if current.Sub(currentWinStart) == swl.windowSize {
			prevWinCount = swl.currWindow.count
		}

		swl.prevWindow.set(currentWinStart, prevWinCount)
		swl.currWindow.set(current, 0)
	}

}
