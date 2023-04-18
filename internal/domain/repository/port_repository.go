package repository

import (
	"context"

	"github.com/canbo-x/port-service/internal/domain/model"
)

// PortRepository defines the interface for the port repository.
type PortRepository interface {
	// Upsert inserts or updates a port in the repository.
	Upsert(ctx context.Context, port *model.Port) error

	// Get returns the port with the given id.
	Get(ctx context.Context, id string) (*model.Port, error)

	// GetLength returns the number of ports in the repository.
	GetLength(ctx context.Context) int
}
