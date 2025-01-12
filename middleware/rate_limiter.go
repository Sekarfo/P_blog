package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter  *rate.Limiter
	mu       sync.Mutex
	lastSeen map[string]time.Time
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiter:  rate.NewLimiter(r, b),
		lastSeen: make(map[string]time.Time),
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		ip := r.RemoteAddr
		now := time.Now()

		// Clean up old entries
		for k, v := range rl.lastSeen {
			if now.Sub(v) > time.Minute {
				delete(rl.lastSeen, k)
			}
		}

		rl.lastSeen[ip] = now

		if !rl.limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
