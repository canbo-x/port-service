// Package service contains the business logic for handling port-related operations.
package service

import (
	"context"

	"github.com/canbo-x/port-service/internal/domain/model"
	"github.com/canbo-x/port-service/internal/domain/repository"
	errs "github.com/canbo-x/port-service/internal/error"
)

// PortService encapsulates the logic for working with ports.
type PortService struct {
	portRepo repository.PortRepository
}

// NewPortService creates a new PortService instance with the given port repository.
func NewPortService(portRepo repository.PortRepository) *PortService {
	return &PortService{
		portRepo: portRepo,
	}
}

// UpsertPort inserts or updates a port in the repository.
// If the port is nil, it returns an ErrInvalidInput error.
func (s *PortService) UpsertPort(ctx context.Context, port *model.Port) error {
	if port == nil {
		return errs.ErrInvalidInput
	}

	return s.portRepo.Upsert(ctx, port)
}

// GetPort retrieves a port from the repository using the provided ID.
// If the ID is invalid, it returns an appropriate error.
// If the port is not found, it returns an ErrPortNotFound error.
func (s *PortService) GetPort(ctx context.Context, id string) (*model.Port, error) {
	if err := model.ValidatePortID(id); err != nil {
		return nil, err
	}

	port, err := s.portRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if port == nil {
		return nil, errs.ErrPortNotFound
	}

	return port, nil
}

// GetLength returns the number of ports stored in the repository.
func (s *PortService) GetLength(ctx context.Context) int {
	return s.portRepo.GetLength(ctx)
}
