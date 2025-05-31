package rate_limiter

import (
	"testing"
	"time"
)

func TestFWLimiter_Allow(t *testing.T) {
	limiter := NewFWLimiter(10, 10*time.Second)

	for i := 1; i <= 10; i++ {
		if !limiter.Allow() && limiter.requests < 10 {
			t.Errorf("10 requests should be allowed but %d requests were allowed", i)
		}
	}
}

func TestFWLimiter_Blocked(t *testing.T) {
	limiter := NewFWLimiter(10, 10*time.Second)

	for i := 1; i <= 10; i++ {
		limiter.Allow()
	}

	if limiter.Allow() && limiter.requests > 10 {
		t.Errorf("Only 10 requests should be allowed but %d requests were allowed", limiter.requests)
	}
}
