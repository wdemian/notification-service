package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/thatva/notification-service/internal/config"
	"github.com/thatva/notification-service/notifier"
	"github.com/thatva/notification-service/ratelimiter"
)

func main() {
	configFile := flag.String("c", "", "Path to the configuration file")
	numMessages := flag.Int("n", 3, "Number of messages to send")
	notificationType := flag.String("t", "news", "Notification Type")
	userList := flag.String("u", "", "Comma-separated list of destination users")
	flag.Parse()

	cfg, err := config.LoadRateLimiterConfig(*configFile)
	if err != nil {
		log.Fatal("error reading configuration file", err)
	}

	rateLimiter := ratelimiter.NewRateLimiter(cfg)
	service := notifier.NewNotifier(notifier.ConsoleGateway{}, rateLimiter)

	users := strings.Split(*userList, ",")
	if len(users) == 0 || users[0] == "" {
		// Send message to User A and B if no users were provided
		users = []string{"User A", "User B"}
	}
	for i := 0; i < *numMessages; i++ {
		for _, u := range users {
			send(service, *notificationType, u, fmt.Sprintf("update %d", i+1))
		}
	}
}

func send(s notifier.NotificationService, notificationType, userID, message string) {
	err := s.Send(notificationType, userID, message)
	if err != nil {
		log.Println(err)
	}
}
