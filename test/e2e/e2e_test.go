package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/canbo-x/port-service/internal/application/filereader"
	"github.com/canbo-x/port-service/internal/application/service"
	"github.com/canbo-x/port-service/internal/domain/model"
	"github.com/canbo-x/port-service/internal/infrastructure/httpserver"
	"github.com/canbo-x/port-service/internal/infrastructure/repository/memory"
	"github.com/canbo-x/port-service/internal/util"
)

func TestE2E(t *testing.T) {
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

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		require.FailNow(t, "failed to get current working directory", err.Error())
	}
	// Get the parent directory of the current working directory
	parentDir := filepath.Dir(cwd)

	// Construct the path to the test data file using filepath.Join()
	filename := filepath.Join(parentDir, "testdata", "ports.json")

	// Initialize the file reader
	fileReader := &filereader.JSONFileReader{
		Filename:   filename,
		BufferSize: 1024,
	}

	wg := &sync.WaitGroup{}

	// Start the file processing
	wg.Add(1)
	go func() {
		err = portService.StoreFileToDB(ctx, fileReader, wg)
		require.NoError(t, err)
	}()

	// Wait for the processing to be finished
	wg.Wait()

	// Check if the context was canceled during file processing
	// This will prevent the server from starting if the file processing was canceled

	if ctx.Err() != nil {
		require.FailNow(t, "File processing was canceled", ctx.Err().Error())
	}

	// Check if the length of the repository is the expected one
	require.Equal(t, 2, portService.GetLength(ctx))

	log.Println("File processing complete. Starting HTTP server.")

	// Initialize the HTTP server
	httpServer := httpserver.NewHTTPServer(portService)

	// Start the server
	wg.Add(1)
	go func() {
		err = httpServer.StartServer(ctx, cancel, wg)
		require.NoError(t, err)
	}()

	// Wait for the server to start
	wg.Wait()

	// Check if the context was canceled during server startup
	if ctx.Err() != nil {
		log.Println("Server startup was canceled")
		require.FailNow(t, "Server startup was canceled", ctx.Err().Error())
	}

	// Check if the server is running using the health check endpoint
	resp, err := http.Get("http://localhost:8080/healthz")
	if err != nil {
		require.FailNow(t, "failed to get health check", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		require.FailNow(t, fmt.Sprintf("expected status code to be %d, got %d", http.StatusOK, resp.StatusCode))
	}

	testCases := []struct {
		name             string
		url              string
		method           string
		expectedStatus   int
		validateResponse func(*testing.T, *http.Response)
	}{
		{
			name:           "test get valid port",
			url:            "http://localhost:8080/ports/GBLON",
			method:         "GET",
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, resp *http.Response) {
				var respPort model.Port
				if err := json.NewDecoder(resp.Body).Decode(&respPort); err != nil {
					require.FailNow(t, "failed to unmarshal response", err.Error())
				}
				expectedPort := getGBLON()
				require.Equal(t, *expectedPort, respPort)
			},
		},
		{
			name:           "test get invalid port",
			url:            "http://localhost:8080/ports/invalid",
			method:         "GET",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "test unsupported method",
			url:            "http://localhost:8080/ports/GBLON",
			method:         "POST",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "test malformed URL",
			url:            "http://localhost:8080/ports/GBLON/some_invalid_path",
			method:         "GET",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		// this will ensure that the test will wait for all the test cases to finish
		// so the server won't be shut down before the test cases are finished
		wg.Add(1)
		go func(tc struct {
			name             string
			url              string
			method           string
			expectedStatus   int
			validateResponse func(*testing.T, *http.Response)
		},
		) {
			defer wg.Done()

			var resp *http.Response
			var err error

			switch tc.method {
			case http.MethodGet:
				resp, err = http.Get(tc.url)
			case http.MethodPost:
				resp, err = http.Post(tc.url, "application/json", nil) // Adjust the content type if needed
			default:
				require.FailNow(t, "unsupported method", tc.method)
			}

			require.NoError(t, err)

			defer resp.Body.Close()

			require.Equalf(t, tc.expectedStatus, resp.StatusCode, "test case: %s", tc.name)
			if tc.validateResponse != nil {
				tc.validateResponse(t, resp)
			}
		}(tc)
	}

	wg.Wait()
}

func getGBLON() *model.Port {
	return &model.Port{
		ID:          "GBLON",
		Name:        "London",
		City:        "London",
		Country:     "United Kingdom",
		Alias:       []string{},
		Regions:     []string{},
		Coordinates: []float64{-0.0833, 51.5},
		Province:    "Greater London",
		Timezone:    "Europe/London",
		Unlocs:      []string{"GBLON"},
		Code:        "12345",
	}
}
