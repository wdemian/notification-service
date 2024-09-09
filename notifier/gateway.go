package notifier

import "fmt"

// Gateway represents the underlying notification delivery mechanism
type Gateway interface {
	Send(userID, message string) error
}

// ConsoleGateway is a simple Gateway implementation that prints to console
type ConsoleGateway struct{}

// Send implements the Gateway interface for ConsoleGateway
func (gateway ConsoleGateway) Send(userID, message string) error {
	fmt.Printf("sent message to user %s: %s\n", userID, message)
	return nil
}
