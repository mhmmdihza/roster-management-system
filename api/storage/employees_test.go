package storage

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestTransactionRollback(t *testing.T) {
	_withTestDatabase(t, func(st *Storage) {
		t.Run("Create 3 records with 1 rollback", func(t *testing.T) {
			ctx := context.Background()

			// 1st record: commit
			txCtx1, err := st.NewTransacton(ctx)
			require.NoError(t, err)

			id1, err := st.CreateNewEmployee(txCtx1, "Employee 1", "ACTIVE", 1)
			require.NoError(t, err)

			err = st.Commit(txCtx1)
			require.NoError(t, err)

			// 2nd record: rollback
			txCtx2, err := st.NewTransacton(ctx)
			require.NoError(t, err)

			id2, err := st.CreateNewEmployee(txCtx2, "Employee 2", "ACTIVE", 1)
			require.NoError(t, err)

			err = st.Rollback(txCtx2)
			require.NoError(t, err)

			// 3rd record: non-transactional
			id3, err := st.CreateNewEmployee(ctx, "Employee 3", "ACTIVE", 1)
			require.NoError(t, err)

			// Validate ID sequences
			assert.Equal(t, id1+1, id2, "ID2 should be next after ID1, even if rolled back")
			assert.Equal(t, id2+1, id3, "ID3 should be next after ID2, regardless of rollback")

			// Check if record with id1 exists
			rec1, err := st.SelectEmployeeByID(ctx, id1)
			assert.NoError(t, err)
			assert.Equal(t, "Employee 1", rec1.Name)

			// Check if record with id2 DOES NOT exist
			_, err = st.SelectEmployeeByID(ctx, id2)
			assert.True(t, errors.Is(err, sql.ErrNoRows))

			// Check if record with id3 exists
			rec3, err := st.SelectEmployeeByID(ctx, id3)
			assert.NoError(t, err)
			assert.Equal(t, "Employee 3", rec3.Name)
		})
	})
}
