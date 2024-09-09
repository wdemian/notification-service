package notifier

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/thatva/notification-service/ratelimiter"
)

type MockGateway struct {
	mu   sync.Mutex
	sent []string
}

func (m *MockGateway) Send(userID, message string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sent = append(m.sent, fmt.Sprintf("user: %s, message: %s", userID, message))
	return nil
}

func (m *MockGateway) SentMessages() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sent
}

func TestNotifier_Send(t *testing.T) {
	config := ratelimiter.RateLimiterConfig{
		Limits: map[string]ratelimiter.LimitConfig{
			"news":   {Limit: 1, Period: 24 * time.Hour},
			"update": {Limit: 2, Period: time.Minute},
		},
	}

	testCases := []struct {
		name             string
		notificationType string
		userID           string
		numMessages      int
		wantSent         int
	}{
		{
			name:             "News: allow only 1 message",
			notificationType: "news",
			userID:           "user1",
			numMessages:      2,
			wantSent:         1,
		},
		{
			name:             "Update: allow 2 messages",
			notificationType: "update",
			userID:           "user2",
			numMessages:      3,
			wantSent:         2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockGateway := &MockGateway{}
			rateLimiter := ratelimiter.NewRateLimiter(config)
			notifier := NewNotifier(mockGateway, rateLimiter)

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()
				for i := 0; i < tc.numMessages; i++ {
					_ = notifier.Send(tc.notificationType, tc.userID, "test message")
				}
			}()

			wg.Wait()

			// Check if the number of messages sent matches the expected value
			sentMessages := mockGateway.SentMessages()
			if len(sentMessages) != tc.wantSent {
				t.Errorf("expected %d messages to be sent, but got %d", tc.wantSent, len(sentMessages))
			}
		})
	}
}
