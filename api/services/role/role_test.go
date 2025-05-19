package role

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"payd/storage"

	"github.com/stretchr/testify/assert"
)

// mockStorage implements Storage interface for testing.
type mockStorage struct {
	rolesToReturn []storage.Role
	errToReturn   error
	mu            sync.Mutex
	callCount     int
}

func (m *mockStorage) SelectAllRoles(ctx context.Context) ([]storage.Role, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	return m.rolesToReturn, m.errToReturn
}

func TestRoleManager(t *testing.T) {
	t.Run("refresh", func(t *testing.T) {
		ctx := context.Background()
		mockSt := &mockStorage{
			rolesToReturn: []storage.Role{
				{ID: 1, Name: "a"},
				{ID: 2, Name: "waiter"},
			},
		}

		rm := NewRoleManager(mockSt, time.Minute)

		// Call refresh and expect roles to be updated
		rm.refresh(ctx)

		gotRoles := rm.GetRoles()
		assert.Equal(t, 2, len(gotRoles))
		assert.Equal(t, "a", gotRoles[0].Name)
		assert.Equal(t, "waiter", gotRoles[1].Name)
	})

}

func TestRoleManager_refreshWithError(t *testing.T) {
	t.Run("error on refresh should give latest role records", func(t *testing.T) {
		ctx := context.Background()

		mockSt := &mockStorage{
			rolesToReturn: []storage.Role{
				{ID: 1, Name: "a"},
				{ID: 2, Name: "waiter"},
			},
		}
		rm := NewRoleManager(mockSt, time.Minute)
		// 1st call succeed
		err := rm.refresh(ctx)
		assert.NoError(t, err)

		gotRoles := rm.GetRoles()
		assert.Equal(t, 2, len(gotRoles))
		assert.Equal(t, "a", gotRoles[0].Name)
		assert.Equal(t, "waiter", gotRoles[1].Name)

		rm.storage = &mockStorage{
			errToReturn: errors.New("db error"),
		}

		// 2nd call failed , but still keeping the latest role records
		err = rm.refresh(ctx)
		assert.ErrorContains(t, err, "db error")

		gotRoles = rm.GetRoles()
		assert.Equal(t, 2, len(gotRoles))
		assert.Equal(t, "a", gotRoles[0].Name)
		assert.Equal(t, "waiter", gotRoles[1].Name)
	})
}

func TestRoleManager_Start(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockSt := &mockStorage{
		rolesToReturn: []storage.Role{
			{ID: 10, Name: "testrole"},
		},
	}

	rm := NewRoleManager(mockSt, time.Second*1)
	rm.Start(ctx)

	gotRoles := rm.GetRoles()
	assert.Len(t, gotRoles, 1)
	assert.Equal(t, "testrole", gotRoles[0].Name)
	mockSt.rolesToReturn = []storage.Role{
		{ID: 1, Name: "a"},
		{ID: 2, Name: "waiter"},
	}

	// add delay to let background updating in-memory records
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	<-ticker.C

	gotRoles = rm.GetRoles()
	assert.Equal(t, 2, len(gotRoles))
	assert.Equal(t, "a", gotRoles[0].Name)
	assert.Equal(t, "waiter", gotRoles[1].Name)
	assert.Greater(t, mockSt.callCount, 1)
}
