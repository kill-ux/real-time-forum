package middlewares

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"forum/utils"
)

// RateLimiter holds the state for rate limiting.
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]time.Time // Map of client IP to last access time
}

// NewRateLimiter creates a new RateLimiter instance.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		clients: map[string]time.Time{},
	}
}

// Allow checks if a request from the given IP is allowed.
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Check if the IP exists in the map
	_, exists := rl.clients[ip]
	if exists {
		// Rate limit the user if they accessed within the last second
		return false
	}

	// Update the last access time
	rl.clients[ip] = time.Now()
	return true
}

// CleanupOldEntries removes entries older than 1 second.
func (rl *RateLimiter) CleanupOldEntries() {
	for {
		rl.mu.Lock()
		for ip, lastAccess := range rl.clients {
			if time.Since(lastAccess) > time.Millisecond*100 {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
		time.Sleep(time.Millisecond * 100) // Run cleanup every 100 ms
	}
}

// RateLimitMiddleware is the HTTP middleware for rate limiting.
func RateLimit(rl *RateLimiter, next http.Handler) http.Handler {
	return ErrorHandler(LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the client's IP address
		ip := strings.Split(r.RemoteAddr, ":")[0] + r.URL.Path

		// Check if the request is allowed
		if !rl.Allow(ip) {
			utils.RespondWithError(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})))
}
