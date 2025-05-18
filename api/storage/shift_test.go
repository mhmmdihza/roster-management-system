package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateShift(t *testing.T) {
	_withTestDatabase(t, func(st *Storage) {
		ctx := context.Background()
		t.Run("Insert and list shifts by time range", func(t *testing.T) {
			dayA := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC)
			shift1Start := dayA.Add(9 * time.Hour)
			shift1End := shift1Start.Add(8 * time.Hour)

			id1, err := st.CreateNewShiftSchedule(ctx, 1, shift1Start, shift1End)
			assert.NoError(t, err)
			assert.Greater(t, id1, 0)

		})

		t.Run("Insert shift with invalid role_id should fail", func(t *testing.T) {
			start := time.Now()
			end := start.Add(8 * time.Hour)
			_, err := st.CreateNewShiftSchedule(ctx, -1, start, end) // role_id 0 is invalid (no FK)
			assert.Error(t, err)
		})
	})
}

func TestGetAvailableShiftsByTimeRangeAndRole(t *testing.T) {
	_withTestDatabase(t, func(st *Storage) {
		ctx := context.Background()
		roleID := 1

		// Setup shifts (some within the range, some outside)
		shifts := []struct {
			role  int
			start time.Time
			end   time.Time
		}{
			{roleID, time.Date(2025, 5, 15, 9, 0, 0, 0, time.UTC), time.Date(2025, 5, 15, 17, 0, 0, 0, time.UTC)},
			{roleID, time.Date(2025, 5, 15, 18, 0, 0, 0, time.UTC), time.Date(2025, 5, 15, 22, 0, 0, 0, time.UTC)},
			{2, time.Date(2025, 5, 16, 9, 0, 0, 0, time.UTC), time.Date(2025, 5, 16, 17, 0, 0, 0, time.UTC)}, // different role
		}

		shiftIDs := make([]int, 0, len(shifts))
		for _, s := range shifts {
			id, err := st.CreateNewShiftSchedule(ctx, s.role, s.start, s.end)
			assert.NoError(t, err)
			assert.Greater(t, id, 0)
			shiftIDs = append(shiftIDs, id)
		}

		// Create employees for shift requests
		employeeID, err := st.CreateNewEmployee(ctx, "Test Emp", "ACTIVE", roleID)
		assert.NoError(t, err)
		assert.Greater(t, employeeID, 0)

		// Create a shift request with APPROVED status for the first shift
		reqID, err := st.CreateShiftRequest(ctx, employeeID, shiftIDs[0])
		assert.NoError(t, err)
		assert.Greater(t, reqID, 0)

		err = st.UpdateShiftRequestStatusByShiftID(ctx, shiftIDs[0], "APPROVED")
		assert.NoError(t, err)

		start := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC)
		end := time.Date(2025, 5, 15, 23, 59, 59, 0, time.UTC)

		t.Run("returns error if start time is zero", func(t *testing.T) {
			_, err := st.GetAvailableShiftsByTimeRangeAndRole(ctx, time.Time{}, end, roleID)
			assert.Error(t, err)
		})

		t.Run("returns error if end time is zero", func(t *testing.T) {
			_, err := st.GetAvailableShiftsByTimeRangeAndRole(ctx, start, time.Time{}, roleID)
			assert.Error(t, err)
		})

		t.Run("returns available shifts excluding approved shift requests", func(t *testing.T) {
			shifts, err := st.GetAvailableShiftsByTimeRangeAndRole(ctx, start, end, roleID)
			assert.NoError(t, err)

			// The first shift is approved and should be excluded, so only one shift should be returned
			assert.Len(t, shifts, 1)

			// Check the returned shift is the second shift for the role
			assert.Equal(t, shiftIDs[1], shifts[0].ID)
			assert.Equal(t, roleID, shifts[0].RoleID)
		})

		t.Run("returns empty slice if no shifts match role", func(t *testing.T) {
			shifts, err := st.GetAvailableShiftsByTimeRangeAndRole(ctx, start, end, 9999) // non-existing role
			assert.NoError(t, err)
			assert.Empty(t, shifts)
		})
	})
}

func TestDeleteShiftByID(t *testing.T) {
	_withTestDatabase(t, func(st *Storage) {
		ctx := context.Background()
		// Setup a role and employee
		roleID := 1
		employeeID, err := st.CreateNewEmployee(ctx, "Tester", "ACTIVE", roleID)
		assert.NoError(t, err)

		// Create a shift
		start := time.Date(2025, 7, 15, 9, 0, 0, 0, time.UTC)
		end := time.Date(2025, 7, 15, 17, 0, 0, 0, time.UTC)

		shiftID, err := st.CreateNewShiftSchedule(ctx, roleID, start, end)
		assert.NoError(t, err)
		assert.Greater(t, shiftID, 0)

		// Create a shift request for that shift
		shiftReqID, err := st.CreateShiftRequest(ctx, employeeID, shiftID)
		assert.NoError(t, err)
		assert.Greater(t, shiftReqID, 0)

		// Delete the shift
		err = st.DeleteShiftById(ctx, shiftID)
		assert.NoError(t, err)

		// Verify shift is deleted using GetAvailableShiftsByTimeRangeAndRole
		availableShifts, err := st.GetAvailableShiftsByTimeRangeAndRole(ctx, start.Add(-time.Hour), end.Add(time.Hour), roleID)
		assert.NoError(t, err)
		for _, s := range availableShifts {
			assert.NotEqual(t, shiftID, s.ID, "Deleted shift should not be in available shifts")
		}

		// Verify shift request is deleted using ListShiftRequestsByFilterAndTimeRange
		filter := ListShiftRequestFilter{
			EmployeeID: employeeID,
			ShiftID:    shiftID,
			RoleID:     roleID,
		}
		requests, err := st.ListShiftRequestsByFilterAndTimeRange(ctx, filter, start.Add(-time.Hour), end.Add(time.Hour))
		assert.NoError(t, err)
		assert.Empty(t, requests, "Shift request for deleted shift should not exist")
	})
}
