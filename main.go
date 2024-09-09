package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/thatva/notification-service/notifier"
	"github.com/thatva/notification-service/ratelimiter"
)

func main() {
	configFile := flag.String("c", "", "Path to the configuration file")
	numMessages := flag.Int("n", 3, "Number of messages to send")
	notificationType := flag.String("t", "news", "Notification Type")
	flag.Parse()

	config, err := loadRateLimiterConfig(*configFile)
	if err != nil {
		log.Fatal("error reading configuration file", err)
	}

	rateLimiter := ratelimiter.NewRateLimiter(config)
	service := notifier.NewNotifier(notifier.ConsoleGateway{}, rateLimiter)

	for i := 0; i < *numMessages; i++ {
		send(service, *notificationType, "user A", fmt.Sprintf("update %d", i+1))
	}
}

func send(s notifier.NotificationService, notificationType, userID, message string) {
	err := s.Send(notificationType, userID, message)
	if err != nil {
		log.Println(err)
	}
}
