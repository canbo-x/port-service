package main

import (
	"context"
	"log"
	"sync"

	"github.com/canbo-x/port-service/internal/application/filereader"
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

	// Initialize the file reader
	fileReader := &filereader.JSONFileReader{
		// this should be a config value
		Filename:   "ports.json",
		BufferSize: 1024,
	}

	wg := &sync.WaitGroup{}

	// Initialize the HTTP server
	httpServer := httpserver.NewHTTPServer(portService)

	// Start the file processing
	wg.Add(1)
	go func() {
		if err := portService.StoreFileToDB(ctx, fileReader, wg); err != nil {
			log.Printf("Error while processing file: %v", err)
			cancel()
		}
	}()
	// Wait for the processing to be finished
	wg.Wait()

	// Check if the context was canceled during file processing
	// This will prevent the server from starting if the file processing was canceled
	if ctx.Err() != nil {
		log.Println("File processing was canceled")
		return
	}

	log.Println("File processing complete. Now starting the HTTP server.")

	// Start the server
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
