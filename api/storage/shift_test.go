package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateShift(t *testing.T) {
	_withTestDatabase(t, func(st *Storage) {
		t.Run("Insert and list shifts by time range", func(t *testing.T) {
			dayA := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC)
			shift1Start := dayA.Add(9 * time.Hour)
			shift1End := shift1Start.Add(8 * time.Hour)

			id1, err := st.CreateNewShiftSchedule(1, shift1Start, shift1End)
			assert.NoError(t, err)
			assert.Greater(t, id1, 0)

		})

		t.Run("Insert shift with invalid role_id should fail", func(t *testing.T) {
			start := time.Now()
			end := start.Add(8 * time.Hour)
			_, err := st.CreateNewShiftSchedule(0, start, end) // role_id 0 is invalid (no FK)
			assert.Error(t, err)
		})
	})
}
