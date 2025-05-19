package role

import (
	"context"
	"log"
	"sync"
	"time"

	"payd/storage"
)

type Role struct {
	ID   int
	Name string
}

type Storage interface {
	SelectAllRoles(ctx context.Context) ([]storage.Role, error)
}

type RoleManagerInterface interface {
	GetRoles() []Role
}

type RoleManager struct {
	storage Storage

	mu    sync.RWMutex
	roles []Role
	tick  time.Duration
}

func NewRoleManager(st Storage, tick time.Duration) *RoleManager {
	return &RoleManager{
		storage: st,
		roles:   make([]Role, 0),
		tick:    tick,
	}
}

// Start begins the periodic refresh of roles and performs an initial fetch.
func (rm *RoleManager) Start(ctx context.Context) error {
	if err := rm.refresh(ctx); err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(rm.tick)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := rm.refresh(ctx); err != nil {
					log.Printf("periodic role fetch failed: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// refresh fetches and updates the role list from storage.
func (rm *RoleManager) refresh(ctx context.Context) error {
	storageRoles, err := rm.storage.SelectAllRoles(ctx)
	if err != nil {
		return err
	}

	newRoles := make([]Role, len(storageRoles))
	for i, sr := range storageRoles {
		newRoles[i] = Role{ID: sr.ID, Name: sr.Name}
	}

	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.roles = newRoles
	return nil
}

// GetRoles returns a thread-safe copy of the latest role list.
func (rm *RoleManager) GetRoles() []Role {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	copied := make([]Role, len(rm.roles))
	copy(copied, rm.roles)
	return copied
}
