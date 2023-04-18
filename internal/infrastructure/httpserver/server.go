// Package httpserver contains the implementation of the HTTP server for the port service.
package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/canbo-x/port-service/internal/application/handler"
	"github.com/canbo-x/port-service/internal/application/service"
)

// HTTPServer represents the main structure for the HTTP server.
type HTTPServer struct {
	portService *service.PortService
}

// NewHTTPServer creates a new instance of HTTPServer with the given port service.
func NewHTTPServer(portService *service.PortService) *HTTPServer {
	return &HTTPServer{
		portService: portService,
	}
}

// StartServer starts the HTTP server, sets up routes and middleware,
// and handles graceful shutdown when the context is canceled or an error occurs.
func (s *HTTPServer) StartServer(
	ctx context.Context,
	cancel context.CancelFunc,
	wg *sync.WaitGroup,
) error {
	// Check if the port service is nil just in case
	if s.portService == nil {
		return fmt.Errorf("port service is nil")
	}
	// Initialize the HTTP server and the port handler
	e := echo.New()
	portHandler := handler.NewPortHandler(s.portService)

	defer func() {
		// Shutdown the HTTP server with a timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := e.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server forced to shutdown: %v", err)
		}
	}()

	// Set the timeouts for the server
	e.Server.ReadTimeout = 10 * time.Second
	e.Server.WriteTimeout = 10 * time.Second

	// Set the idle timeout for the server
	e.Server.IdleTimeout = 30 * time.Second

	// Set the maximum header size
	e.Server.MaxHeaderBytes = 100 * 1024 // 100 KB

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet},
	}))

	// This is just for demonstration purposes
	// In production, more sophisticated rate limiting should be used
	// Limit the number of requests to 10 per second
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))

	// Add health check endpoint
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	// Routes
	e.GET("/ports/:id", portHandler.GetPort)

	// Start the HTTP server
	serverErrors := make(chan error)
	go func() {
		// Port should be configurable and not hard-coded
		// This is just for demonstration purposes
		// Configuration file logic is not implemented
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Signal that the server has started
	wg.Done()

	// Wait for the context to be canceled or for the server to fail
	select {
	case <-ctx.Done():
		// Context was canceled
		log.Println("Context canceled")
		return ctx.Err()

	case err := <-serverErrors:
		// HTTP server encountered an error
		log.Printf("HTTP server shutting down: %v", err)
		return err
	}
}
