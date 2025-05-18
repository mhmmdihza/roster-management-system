package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectAllRoles(t *testing.T) {
	_withTestDatabase(t, func(st *Storage) {
		ctx := context.Background()

		t.Run("Should return all predefined roles", func(t *testing.T) {
			roles, err := st.SelectAllRoles(ctx)
			require.NoError(t, err)
			require.Len(t, roles, 4)

			expected := map[int]string{
				0: "A",
				1: "Cashier",
				2: "Cook",
				3: "Waiter",
			}

			for _, role := range roles {
				name, ok := expected[role.ID]
				assert.True(t, ok, "Unexpected role ID: %d", role.ID)
				assert.Equal(t, name, role.Name)
			}
		})
	})
}
