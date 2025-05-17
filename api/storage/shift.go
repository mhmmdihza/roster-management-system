package storage

import (
	"context"
	"fmt"
	"time"
)

type Shift struct {
	ID        int       `db:"id"`
	RoleID    int       `db:"role_id"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
	CreatedAt time.Time `db:"created_at"`
}

func (s *Storage) CreateNewShiftSchedule(ctx context.Context, roleId int, startTime, endTime time.Time) (int, error) {
	var id int
	query := `INSERT INTO shifts (role_id, start_time, end_time) VALUES ($1, $2, $3) RETURNING id`
	err := s.db.QueryRowContext(ctx, query, roleId, startTime, endTime).Scan(&id)
	return id, err
}

func (s *Storage) GetAvailableShiftsByTimeRangeAndRole(ctx context.Context, start, end time.Time, roleId int) ([]Shift, error) {
	if start.IsZero() || end.IsZero() {
		return nil, fmt.Errorf("both start and end time must be provided")
	}
	var shifts []Shift
	query := `
        SELECT s.*
        FROM shifts s
        WHERE s.role_id = $3
          AND s.start_time >= $1
          AND s.start_time <= $2
          AND s.id NOT IN (
              SELECT sr.shift_id
              FROM shift_requests sr
              WHERE sr.status = 'APPROVED'
          )
        ORDER BY s.start_time
    `
	err := s.db.SelectContext(ctx, &shifts, query, start, end, roleId)
	return shifts, err
}

func (s *Storage) DeleteShiftById(ctx context.Context, shiftId int) error {
	query := `DELETE FROM shifts WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, shiftId)
	return err
}
