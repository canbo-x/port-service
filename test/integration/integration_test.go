package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/canbo-x/port-service/internal/application/handler"
	"github.com/canbo-x/port-service/internal/application/service"
	"github.com/canbo-x/port-service/internal/domain/model"
	"github.com/canbo-x/port-service/internal/infrastructure/repository/memory"
)

func TestGetPort(t *testing.T) {
	// Create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the repository and the service
	portRepository := memory.NewMemoryDB()
	portService := service.NewPortService(portRepository)

	// Populate the repository with test data
	if err := portRepository.Upsert(ctx, getGBLON()); err != nil {
		t.Fatalf("failed to upsert port: %v", err)
	}

	if err := portRepository.Upsert(ctx, getFRPAR()); err != nil {
		t.Fatalf("failed to upsert port: %v", err)
	}

	// Create the port handler
	portHandler := handler.NewPortHandler(portService)

	testCases := []struct {
		name           string
		id             string
		expectedStatus int
		expectedPort   *model.Port
	}{
		{
			name:           "Valid Port ID - GBLON",
			id:             "GBLON",
			expectedStatus: http.StatusOK,
			expectedPort:   getGBLON(),
		},
		{
			name:           "Valid Port ID - FRPAR",
			id:             "FRPAR",
			expectedStatus: http.StatusOK,
			expectedPort:   getFRPAR(),
		},
		{
			name:           "Invalid Port ID",
			id:             "invalid",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create a test request
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/ports/%s", tc.id), nil)
			rec := httptest.NewRecorder()

			// Create Echo context and set parameters
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/ports/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.id)

			// Execute the handler
			if err := portHandler.GetPort(c); err != nil {
				t.Errorf("handler error: %v", err)
			}

			// Check the response status code
			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedStatus == http.StatusOK {
				// Decode the response
				var respPort model.Port
				if err := json.Unmarshal(rec.Body.Bytes(), &respPort); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				// Compare the results
				assert.Equal(t, *tc.expectedPort, respPort)
			}
		})
	}
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

func getFRPAR() *model.Port {
	return &model.Port{
		ID:          "FRPAR",
		Name:        "Paris",
		City:        "Paris",
		Country:     "France",
		Alias:       []string{},
		Regions:     []string{},
		Coordinates: []float64{2.3488, 48.8534},
		Province:    "Île-de-France",
		Timezone:    "Europe/Paris",
		Unlocs:      []string{"FRPAR"},
		Code:        "23456",
	}
}
