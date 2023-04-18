package util

import (
	"log"
	"time"
)

// GracefulShutdown demonstrates a graceful shutdown by sleeping for 1 second.
func GracefulShutdown() {
	// Perform any additional cleanup if necessary
	log.Println("Service gracefully shut down requested")
	time.Sleep(1 * time.Second)
	log.Println("Service gracefully shut down completed")
}
