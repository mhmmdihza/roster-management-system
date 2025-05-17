package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInsertAndSelectRecord(t *testing.T) {
	_withTestDatabase(t, func(st *Storage) {
		ctx := context.Background()
		t.Run("Valid insert and select", func(t *testing.T) {
			name := "Test Employee"
			status := "ACTIVE"

			id, err := st.CreateNewEmployee(ctx, name, status, 1)
			assert.NoError(t, err)
			assert.Greater(t, id, 0)

			record, err := st.SelectEmployeeByID(ctx, id)
			assert.NoError(t, err)
			assert.Equal(t, name, record.Name)
			assert.Equal(t, status, record.Status)
			assert.WithinDuration(t, time.Now(), record.CreatedAt, 5*time.Second)
		})

		t.Run("Invalid status", func(t *testing.T) {
			_, err := st.CreateNewEmployee(ctx, "Invalid Status Employee", "UNKNOWN", 1)
			assert.Error(t, err)
		})

		t.Run("Empty name", func(t *testing.T) {
			_, err := st.CreateNewEmployee(ctx, "", "ACTIVE", 1)
			assert.Error(t, err)
		})
		t.Run("Invalid Role", func(t *testing.T) {
			_, err := st.CreateNewEmployee(ctx, "Invalid Role", "ACTIVE", 0)
			assert.Error(t, err)
		})
	})
}

func TestUpdateStatus(t *testing.T) {
	_withTestDatabase(t, func(st *Storage) {
		ctx := context.Background()
		t.Run("Update one employee out of two", func(t *testing.T) {
			id1, err := st.CreateNewEmployee(ctx, "Emp One", "ACTIVE", 1)
			assert.NoError(t, err)

			id2, err := st.CreateNewEmployee(ctx, "Emp Two", "ACTIVE", 1)
			assert.NoError(t, err)

			err = st.UpdateEmployeeStatus(ctx, id2, "INACTIVE")
			assert.NoError(t, err)

			// Validate first employee remains unchanged
			emp1, err := st.SelectEmployeeByID(ctx, id1)
			assert.NoError(t, err)
			assert.Equal(t, "ACTIVE", emp1.Status)

			// Validate second employee got updated
			emp2, err := st.SelectEmployeeByID(ctx, id2)
			assert.NoError(t, err)
			assert.Equal(t, "INACTIVE", emp2.Status)
		})
	})
}
