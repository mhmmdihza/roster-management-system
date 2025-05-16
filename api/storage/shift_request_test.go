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
