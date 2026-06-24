package middleware

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(3, time.Second)

	for i := 0; i < 3; i++ {
		if !limiter.Allow("test") {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	if limiter.Allow("test") {
		t.Error("4th request should be denied")
	}
}

func TestRateLimiterWindow(t *testing.T) {
	limiter := NewRateLimiter(2, 10*time.Millisecond)

	limiter.Allow("test")
	limiter.Allow("test")

	time.Sleep(15 * time.Millisecond)

	if !limiter.Allow("test") {
		t.Error("Request after window should be allowed")
	}
}
