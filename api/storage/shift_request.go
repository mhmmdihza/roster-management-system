package storage

import (
	"context"
	"fmt"
	"time"
)

type ShiftRequestWithShiftDetails struct {
	ID           int        `db:"id"`
	EmployeeID   int        `db:"employee_id"`
	EmployeeName string     `db:"employee_name"`
	ShiftID      int        `db:"shift_id"`
	Status       string     `db:"status"`
	RequestedAt  time.Time  `db:"requested_at"`
	ReviewedAt   *time.Time `db:"reviewed_at"`
	ReviewedBy   *int       `db:"reviewed_by"`
	RoleID       int        `db:"role_id"`
	RoleName     string     `db:"role_name"`
	StartTime    time.Time  `db:"start_time"`
	EndTime      time.Time  `db:"end_time"`
}

type ListShiftRequestFilter struct {
	EmployeeID int
	ShiftID    int
	RoleID     int
	Status     string
}

func (s *Storage) CreateShiftRequest(ctx context.Context, employeeId, shiftId int) (int, error) {
	var id int
	query := `
		INSERT INTO shift_requests (employee_id, shift_id, status)
		VALUES ($1, $2, 'PENDING')
		RETURNING id
	`
	err := s.db.QueryRowxContext(ctx, query, employeeId, shiftId).Scan(&id)
	return id, err
}

func (s *Storage) UpdateShiftRequestStatusByShiftID(ctx context.Context, shiftId int, status string) error {
	query := `
		UPDATE shift_requests
		SET status = $1
		WHERE shift_id = $2
	`
	_, err := s.db.ExecContext(ctx, query, status, shiftId)
	return err
}

func (s *Storage) ListShiftRequestsByFilterAndTimeRange(ctx context.Context, filter ListShiftRequestFilter,
	start time.Time,
	end time.Time) ([]ShiftRequestWithShiftDetails, error) {
	if start.IsZero() || end.IsZero() {
		return nil, fmt.Errorf("both start and end time must be provided")
	}
	baseQuery := `
        SELECT 
            sr.id, sr.employee_id, e.name AS employee_name, 
			sr.shift_id, sr.status, sr.requested_at, sr.reviewed_at, sr.reviewed_by,
            s.role_id, r.name AS role_name, s.start_time, s.end_time
        FROM 
            shift_requests sr
        JOIN 
            shifts s ON sr.shift_id = s.id
		JOIN 
			employees e ON sr.employee_id = e.id
		JOIN 
    		roles r ON s.role_id = r.id
        WHERE s.start_time >= $1 AND s.start_time <= $2 
    `

	args := []interface{}{start, end}
	argPos := len(args) + 1

	if filter.EmployeeID != 0 {
		baseQuery += fmt.Sprintf(" AND sr.employee_id = $%d", argPos)
		args = append(args, filter.EmployeeID)
		argPos++
	}

	if filter.ShiftID != 0 {
		baseQuery += fmt.Sprintf(" AND sr.shift_id = $%d", argPos)
		args = append(args, filter.ShiftID)
		argPos++
	}

	if filter.RoleID != 0 {
		baseQuery += fmt.Sprintf(" AND s.role_id = $%d", argPos)
		args = append(args, filter.RoleID)
		argPos++
	}

	if filter.Status != "" {
		baseQuery += fmt.Sprintf(" AND sr.status = $%d", argPos)
		args = append(args, filter.Status)
		argPos++
	}

	baseQuery += " ORDER BY s.start_time ASC"

	var results []ShiftRequestWithShiftDetails
	err := s.db.SelectContext(ctx, &results, baseQuery, args...)
	return results, err
}
