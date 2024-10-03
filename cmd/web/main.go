package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/thatva/notification-service/internal/config"
	"github.com/thatva/notification-service/notifier"
	"github.com/thatva/notification-service/ratelimiter"
)

type Request struct {
	Type    string `json:"type"`
	UserID  string `json:"userid"`
	Mesasge string `json:"message"`
}

func main() {
	configFile := flag.String("c", "", "Path to the configuration file")
	port := flag.Int("p", 8080, "Port to listen on")
	flag.Parse()

	cfg, err := config.LoadRateLimiterConfig(*configFile)
	if err != nil {
		log.Fatal("error reading configuration file", err)
	}

	rateLimiter := ratelimiter.NewRateLimiter(cfg)
	service := notifier.NewNotifier(notifier.ConsoleGateway{}, rateLimiter)

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/notification", sendNotification(service))
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Listening on %s", addr)
	err = http.ListenAndServe(addr, mux)
	log.Fatal(err)
}

func sendNotification(service notifier.NotificationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
		}
		err := service.Send(req.Type, req.UserID, req.Mesasge)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusTooManyRequests)
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello")
}
