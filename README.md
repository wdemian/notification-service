# Notification Service

A (very) basic Go-based test project to send notifications with rate limiting.

## Features

* Sends notifications to users.
* Implements rate limiting to prevent spamming.
* Configurable limits per notification type.
* Pluggable gateway interface for different delivery mechanisms (only a dummy console notification is provided)

## Project Structure

```
├── internal                
    └── config               # Configuration loading
├── main.go                  # Main application entry point
├── cmd
│   ├── web                  # Web service entry point
│   └── cli                  # CLI tool to test notifications
├── notifier
│   ├── gateway.go           # Gateway interface and console implementation
│   ├── notifier.go          # Notification service implementation
│   └── notifier_test.go     # Unit tests for the notifier
└── ratelimiter
    ├── ratelimiter.go       # Rate limiter implementation
    └── ratelimiter_test.go  # Unit tests for the rate limiter
```

## Getting Started

   ```bash
   git clone https://github.com/thatva/notification-service.git
   make deps && make build
   ```

```bash
   bin/notification-service -h
   Usage of bin/notification-service:
  -c string
        Path to the configuration file
  -p int
        Port to listen on (default 8080)

 Example: bin/notification-service -c config.yaml -p 8080
```

The cli tool has some extra options for testing:

```bash
   bin/notification-service-cli -h
   Usage of bin/notification-service-cli:
  -c string
        Path to the configuration file
  -n int
        Number of messages to send (default 3)
  -t string
        Notification Type (default "news")
  -u string
        Comma-separated list of destination users

 Example: bin/notification-service -c config.yaml -n 5 -t update -c userA,userB

```

### Configuration

The service can be configured using a YAML file. See `config.yaml` for an example configuration.
A notification type **must** have configured limits, otherwise it will be rejected
