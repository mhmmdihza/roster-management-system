package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateShiftRequest(t *testing.T) {
	_withTestDatabase(t, func(st *Storage) {
		// Create an employee
		employeeName := "emp 1"
		employeeStatus := "ACTIVE"
		roleID := 1
		employeeId, err := st.CreateNewEmployee(employeeName, employeeStatus, roleID)
		assert.NoError(t, err)
		assert.Greater(t, employeeId, 0)

		// Create a shift
		startTime := time.Now().Add(1 * time.Hour)
		endTime := startTime.Add(8 * time.Hour)
		shiftId, err := st.CreateNewShiftSchedule(roleID, startTime, endTime)
		assert.NoError(t, err)
		assert.Greater(t, shiftId, 0)
		t.Run("Valid shift request insert", func(t *testing.T) {
			requestID, err := st.CreateShiftRequest(employeeId, shiftId)
			assert.NoError(t, err)
			assert.Greater(t, requestID, 0)
		})
		t.Run("Invalid shift id", func(t *testing.T) {
			_, err := st.CreateShiftRequest(employeeId, 100)
			assert.Error(t, err)
		})
		t.Run("Invalid employee id", func(t *testing.T) {
			_, err := st.CreateShiftRequest(100, shiftId)
			assert.Error(t, err)
		})
	})
}

func TestListShiftRequestsByFilterAndTimeRange(t *testing.T) {
	_withTestDatabase(t, func(st *Storage) {
		// Setup employees
		employees := []struct {
			name   string
			status string
			role   int
		}{
			{"Alice", "ACTIVE", 1},
			{"Bob", "ACTIVE", 2},
			{"Charlie", "ACTIVE", 3},
		}

		for _, emp := range employees {
			employeeID, err := st.CreateNewEmployee(emp.name, emp.status, emp.role)
			assert.NoError(t, err)
			assert.Greater(t, employeeID, 0)
		}

		// Setup shifts
		shiftTimes := []struct {
			role  int
			start time.Time
			end   time.Time
		}{
			{1, time.Date(2025, 5, 15, 9, 0, 0, 0, time.UTC), time.Date(2025, 5, 15, 17, 0, 0, 0, time.UTC)},
			{2, time.Date(2025, 5, 15, 13, 0, 0, 0, time.UTC), time.Date(2025, 5, 15, 21, 0, 0, 0, time.UTC)},
			{1, time.Date(2025, 5, 16, 9, 0, 0, 0, time.UTC), time.Date(2025, 5, 16, 17, 0, 0, 0, time.UTC)},
		}
		for _, stf := range shiftTimes {
			idShift, err := st.CreateNewShiftSchedule(stf.role, stf.start, stf.end)
			assert.NoError(t, err)
			assert.Greater(t, idShift, 0)
		}

		// Setup shift requests
		shiftRequest := []struct {
			employeeId int
			shiftId    int
		}{
			{1, 1}, // This will insert with status 'PENDING' by default
			{2, 2},
			{1, 3},
		}

		for _, req := range shiftRequest {
			idShiftRequest, err := st.CreateShiftRequest(req.employeeId, req.shiftId)
			assert.NoError(t, err)
			assert.Greater(t, idShiftRequest, 0)
		}

		start := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC)
		end := time.Date(2025, 5, 16, 23, 59, 59, 0, time.UTC)

		t.Run("return error if start time is zero", func(t *testing.T) {
			_, err := st.ListShiftRequestsByFilterAndTimeRange(ListShiftRequestFilter{}, time.Time{}, end)
			assert.Error(t, err)
		})

		t.Run("return error if end time is zero", func(t *testing.T) {
			_, err := st.ListShiftRequestsByFilterAndTimeRange(ListShiftRequestFilter{}, start, time.Time{})
			assert.Error(t, err)
		})

		t.Run("list all requests within time range ordered by shift start_time", func(t *testing.T) {
			results, err := st.ListShiftRequestsByFilterAndTimeRange(ListShiftRequestFilter{}, start, end)
			assert.NoError(t, err)
			assert.Len(t, results, 3)

			// Check ascending order by shift start_time
			for i := 0; i < len(results)-1; i++ {
				assert.True(t, results[i].StartTime.Before(results[i+1].StartTime) || results[i].StartTime.Equal(results[i+1].StartTime))
			}

			// Check some expected values
			assert.Equal(t, 1, results[0].ID)
			assert.Equal(t, 2, results[1].ID)
			assert.Equal(t, 3, results[2].ID)
		})

		t.Run("filter by EmployeeID", func(t *testing.T) {
			filter := ListShiftRequestFilter{EmployeeID: 1}
			results, err := st.ListShiftRequestsByFilterAndTimeRange(filter, start, end)
			assert.NoError(t, err)
			assert.Len(t, results, 2)
			for _, r := range results {
				assert.Equal(t, 1, r.EmployeeID)
			}
		})

		t.Run("filter by ShiftID", func(t *testing.T) {
			filter := ListShiftRequestFilter{ShiftID: 2}
			results, err := st.ListShiftRequestsByFilterAndTimeRange(filter, start, end)
			assert.NoError(t, err)
			assert.Len(t, results, 1)
			assert.Equal(t, 2, results[0].ShiftID)
		})

		t.Run("filter by RoleID", func(t *testing.T) {
			filter := ListShiftRequestFilter{RoleID: 1}
			results, err := st.ListShiftRequestsByFilterAndTimeRange(filter, start, end)
			assert.NoError(t, err)
			assert.Len(t, results, 2)
		})

		t.Run("filter by status", func(t *testing.T) {
			filter := ListShiftRequestFilter{Status: "PENDING"}
			results, err := st.ListShiftRequestsByFilterAndTimeRange(filter, start, end)
			assert.NoError(t, err)
			assert.Len(t, results, 3)
		})

		t.Run("filter by EmployeeID and ShiftID", func(t *testing.T) {
			filter := ListShiftRequestFilter{EmployeeID: 1, ShiftID: 3}
			results, err := st.ListShiftRequestsByFilterAndTimeRange(filter, start, end)
			assert.NoError(t, err)
			assert.Len(t, results, 1)
			r := results[0]
			assert.Equal(t, 1, r.EmployeeID)
			assert.Equal(t, 3, r.ShiftID)
		})

		t.Run("filter by EmployeeID and RoleID", func(t *testing.T) {
			filter := ListShiftRequestFilter{EmployeeID: 1, RoleID: 1}
			results, err := st.ListShiftRequestsByFilterAndTimeRange(filter, start, end)
			assert.NoError(t, err)
			assert.Len(t, results, 2)
			for _, r := range results {
				assert.Equal(t, 1, r.EmployeeID)
				assert.Equal(t, 1, r.RoleID)
			}
		})

		t.Run("filter by ShiftID and RoleID", func(t *testing.T) {
			filter := ListShiftRequestFilter{ShiftID: 1, RoleID: 1}
			results, err := st.ListShiftRequestsByFilterAndTimeRange(filter, start, end)
			assert.NoError(t, err)
			assert.Len(t, results, 1)
			r := results[0]
			assert.Equal(t, 1, r.ShiftID)
			assert.Equal(t, 1, r.RoleID)
		})

		t.Run("filter by EmployeeID, ShiftID and RoleID", func(t *testing.T) {
			filter := ListShiftRequestFilter{EmployeeID: 1, ShiftID: 1, RoleID: 1}
			results, err := st.ListShiftRequestsByFilterAndTimeRange(filter, start, end)
			assert.NoError(t, err)
			assert.Len(t, results, 1)
			r := results[0]
			assert.Equal(t, 1, r.EmployeeID)
			assert.Equal(t, 1, r.ShiftID)
			assert.Equal(t, 1, r.RoleID)
		})

		t.Run("filter by start_time between given range", func(t *testing.T) {
			partialStart := time.Date(2025, 5, 15, 12, 0, 0, 0, time.UTC)
			partialEnd := time.Date(2025, 5, 15, 23, 59, 59, 0, time.UTC)

			results, err := st.ListShiftRequestsByFilterAndTimeRange(ListShiftRequestFilter{}, partialStart, partialEnd)
			assert.NoError(t, err)

			// only 1 shift has start_time in this range (shift id 2)
			assert.Len(t, results, 1)

			assert.Equal(t, 2, results[0].ShiftID)
		})
	})
}
