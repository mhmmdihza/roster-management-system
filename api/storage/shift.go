package storage

import (
	"time"
)

func (s *Storage) CreateNewShiftSchedule(roleId int, startTime, endTime time.Time) (int, error) {
	var id int
	query := `INSERT INTO shifts (role_id, start_time, end_time) VALUES ($1, $2, $3) RETURNING id`
	err := s.db.QueryRow(query, roleId, startTime, endTime).Scan(&id)
	return id, err
}
