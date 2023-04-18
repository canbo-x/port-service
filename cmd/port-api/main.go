package main

import (
	"context"
	"log"
	"sync"

	"github.com/canbo-x/port-service/internal/application/service"
	"github.com/canbo-x/port-service/internal/infrastructure/httpserver"
	"github.com/canbo-x/port-service/internal/infrastructure/repository/memory"
	"github.com/canbo-x/port-service/internal/util"
)

func main() {
	// Create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up a signal handler to listen for signals
	util.SetupSignalHandler(ctx, cancel)

	// Shutdown the server gracefully
	defer util.GracefulShutdown()

	// Initialize the repository and the service
	portRepository := memory.NewMemoryDB()
	portService := service.NewPortService(portRepository)

	// Initialize the HTTP server
	httpServer := httpserver.NewHTTPServer(portService)

	// Start the server
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := httpServer.StartServer(ctx, cancel, wg); err != nil {
			log.Printf("Error starting HTTP server: %v", err)
			cancel()
		}
	}()

	// Wait for the server to start
	wg.Wait()

	// Wait for the context to be canceled
	<-ctx.Done()
}
