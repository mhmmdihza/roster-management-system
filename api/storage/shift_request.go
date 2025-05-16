package storage

func (s *Storage) CreateShiftRequest(employeeId, shiftId int) (int, error) {
	var id int
	query := `
		INSERT INTO shift_requests (employee_id, shift_id, status)
		VALUES ($1, $2, 'PENDING')
		RETURNING id
	`
	err := s.db.QueryRow(query, employeeId, shiftId).Scan(&id)
	return id, err
}
