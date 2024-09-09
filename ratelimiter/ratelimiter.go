package ratelimiter

import (
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// LimitConfig defines the rate limit configuration for a notification type
type LimitConfig struct {
	Limit  int
	Period time.Duration
}

// RateLimiterConfig holds the configuration for all notification types
type RateLimiterConfig struct {
	Limits map[string]LimitConfig
}

// RateLimiter manages rate limits for various keys
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	config   RateLimiterConfig
	mu       sync.Mutex
}

type NotificationRequest struct {
	UserID           string
	NotificationType string
}

func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		config:   config,
		mu:       sync.Mutex{},
	}
}

// Allow checks if a notification request is allowed based on rate limits
func (rl *RateLimiter) Allow(request NotificationRequest) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	key := fmt.Sprintf("%s:%s", request.NotificationType, request.UserID)

	limiter, ok := rl.limiters[key]
	if !ok {
		// Create a new limiter for the key if it doesn't exist
		limitConfig, ok := rl.config.Limits[request.NotificationType]
		if ok {
			limiter = rate.NewLimiter(rate.Every(limitConfig.Period), limitConfig.Limit)
		} else {
			// deny the request if NotificationType has no limit set
			log.Printf("Unknown notification type: %s", request.NotificationType)
			return false
		}
		rl.limiters[key] = limiter
	}

	return limiter.Allow()
}
