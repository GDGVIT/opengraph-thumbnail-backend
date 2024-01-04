package main

import (
	"github.com/GDGVIT/opengraph-thumbnail-backend/cmd"
)

// Message represents the message structure you expect to send to the RabbitMQ queue.

func main() {
	cmd.Execute()
}
