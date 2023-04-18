package memory

import (
	"context"
	"sync"

	"github.com/canbo-x/port-service/internal/domain/model"
	"github.com/canbo-x/port-service/internal/domain/repository"
)

// MemoryDB represents an in-memory database for ports.
type MemoryDB struct {
	mu    sync.RWMutex
	ports map[string]*model.Port
}

// NewMemoryDB creates a new instance of MemoryDB.
func NewMemoryDB() repository.PortRepository {
	return &MemoryDB{
		ports: make(map[string]*model.Port),
	}
}

// Upsert inserts or updates a port in the memory database.
// Please read the readme file for more information about the context.
// This is just a demonstration and more details can be found in the `Personal Thoughts and Notes` section.
func (db *MemoryDB) Upsert(ctx context.Context, port *model.Port) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		db.ports[port.ID] = port
	}

	return nil
}

// Get returns a port from the memory database.
func (db *MemoryDB) Get(ctx context.Context, id string) (*model.Port, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		port, ok := db.ports[id]
		if !ok {
			// Discussion: errs.ErrPortNotFound vs nil
			return nil, nil
		}

		return port, nil
	}
}

// GetLength returns the number of ports in the memory database.
// no test provided for this method to not complicate the example
func (db *MemoryDB) GetLength(ctx context.Context) int {
	db.mu.RLock()
	defer db.mu.RUnlock()

	select {
	case <-ctx.Done():
		return 0
	default:
		return len(db.ports)
	}
}
