package memory

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/canbo-x/port-service/internal/domain/model"
	"github.com/canbo-x/port-service/internal/domain/repository"
)

func createPort() *model.Port {
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

func TestMemoryDB(t *testing.T) {
	testCases := []struct {
		name     string
		testFunc func(t *testing.T, db repository.PortRepository)
	}{
		{
			name: "UpsertAndGet",
			testFunc: func(t *testing.T, db repository.PortRepository) {
				port := createPort()

				// Test Upsert
				ctx := context.Background()
				err := db.Upsert(ctx, port)
				require.NoError(t, err)

				// Test Get
				retrievedPort, err := db.Get(ctx, "GBLON")
				require.NoError(t, err)
				assert.Equal(t, port, retrievedPort)
			},
		},
		{
			name: "GetWithNonExistentID",
			testFunc: func(t *testing.T, db repository.PortRepository) {
				// Test Get with non-existent ID
				retrievedPort, err := db.Get(context.Background(), "NON_EXISTENT")
				require.NoError(t, err)
				assert.Nil(t, retrievedPort)
			},
		},
		{
			name: "ContextCancellation_Upsert",
			testFunc: func(t *testing.T, db repository.PortRepository) {
				port := createPort()

				// Test context cancellation for Upsert
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
				defer cancel()
				time.Sleep(2 * time.Millisecond)
				err := db.Upsert(ctx, port)
				assert.Error(t, err)
			},
		},
	}

	db := NewMemoryDB()
	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.testFunc(t, db)
		})
	}
}
