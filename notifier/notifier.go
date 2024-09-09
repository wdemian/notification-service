package notifier

import (
	"fmt"

	"github.com/thatva/notification-service/ratelimiter"
)

type Notifier struct {
	gateway     Gateway
	rateLimiter *ratelimiter.RateLimiter
}

func NewNotifier(gateway Gateway, rateLimiter *ratelimiter.RateLimiter) NotificationService {
	return &Notifier{gateway: gateway, rateLimiter: rateLimiter}
}

type NotificationService interface {
	Send(notificationType, userID, message string) error
}

func (n *Notifier) Send(notificationType, userID, message string) error {
	request := ratelimiter.NotificationRequest{
		UserID:           userID,
		NotificationType: notificationType,
	}

	if !n.rateLimiter.Allow(request) {
		return fmt.Errorf("rate limit exceeded for user %s and type %s", userID, notificationType)
	}

	msg := fmt.Sprintf("[%s] %s", notificationType, message)
	if err := n.gateway.Send(userID, msg); err != nil {
		return err
	}
	return nil
}
