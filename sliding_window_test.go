package rate_limiter

import (
	"testing"
	"time"
)

const (
	limit = 50
	size  = 60 * time.Second
)

func TestSWLimiter_Allowed(t *testing.T) {
	now := time.Date(2024, 12, 31, 15, 16, 15, 0, time.UTC)
	limiter := NewSWLimiter(limit, size)

	limiter.currWindow.start = now.Truncate(size)
	limiter.prevWindow.start = now.Truncate(size).Add(-size)
	limiter.prevWindow.count = 42
	limiter.currWindow.count = 18

	// approximation should be 49
	if !limiter.process(now) {
		t.Errorf("1 requess should be allowed but none were allowed")
	}
}

func TestSWLimiter_Blocked(t *testing.T) {
	now := time.Date(2024, 12, 31, 15, 16, 5, 0, time.UTC)
	limiter := NewSWLimiter(limit, size)

	limiter.currWindow.start = now.Truncate(size)
	limiter.prevWindow.start = now.Truncate(size).Add(-size)
	limiter.prevWindow.count = 42
	limiter.currWindow.count = 18

	// approximation should be 56
	if limiter.process(now) {
		t.Errorf("No request should be allowed")
	}
}
