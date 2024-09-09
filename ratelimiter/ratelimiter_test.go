package ratelimiter

import (
	"fmt"
	"testing"
	"time"
)

func TestRateLimiterAllow(t *testing.T) {
	config := RateLimiterConfig{
		Limits: map[string]LimitConfig{
			"news":     {Limit: 1, Period: 24 * time.Hour},
			"update":   {Limit: 2, Period: time.Minute},
			"frequent": {Limit: 1, Period: time.Second},
		},
	}

	testCases := []struct {
		name             string
		user             string
		notificationType string
		numMessages      int
		sendDelay        time.Duration
		wantError        bool
	}{
		{
			name:             "news: within limit",
			user:             "A",
			notificationType: "news",
			numMessages:      1,
			wantError:        false,
		},
		{
			name:             "news: exceeds limit",
			user:             "B",
			notificationType: "news",
			numMessages:      2,
			wantError:        true,
		},
		{
			name:             "update: within limit",
			user:             "A",
			notificationType: "update",
			numMessages:      2,
			wantError:        false,
		},
		{
			name:             "update: exceeds limit with delay",
			user:             "B",
			notificationType: "update",
			numMessages:      3,
			sendDelay:        1 * time.Second,
			wantError:        true,
		},
		{
			name:             "frequent: within limit with delay",
			user:             "A",
			notificationType: "frequent",
			numMessages:      3,
			sendDelay:        1 * time.Second,
			wantError:        false,
		},
		{
			name:             "frequent: exceeds limit",
			user:             "B",
			notificationType: "frequent",
			numMessages:      3,
			wantError:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rl := NewRateLimiter(config)

			req := NotificationRequest{UserID: tc.user, NotificationType: tc.notificationType}
			err := sendMessages(t, rl, req, tc.numMessages, tc.sendDelay)
			if (err != nil) != tc.wantError {
				if err != nil {
					t.Error(err)
				} else {
					t.Error("wanted an error but got none")
				}
			}
		})
	}
}

func sendMessages(t *testing.T, rl *RateLimiter, req NotificationRequest, numMessages int, delay time.Duration) error {
	t.Helper()
	for i := 0; i < numMessages; i++ {
		if !rl.Allow(req) {
			return fmt.Errorf("message to %v num %d failed", req, i+1)
		}
		if delay > 0 {
			time.Sleep(delay)
		}
	}
	return nil
}
