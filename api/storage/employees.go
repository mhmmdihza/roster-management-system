package storage

import (
	"context"
	"time"
)

type Employee struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}

func (s *Storage) CreateNewEmployee(ctx context.Context, name, status string, roleId int) (int, error) {
	var id int
	query := `INSERT INTO employees (name, status, role_id) VALUES ($1, $2, $3) RETURNING id`
	err := s.db.QueryRowContext(ctx, query, name, status, roleId).Scan(&id)
	return id, err
}

// UpdateEmployeeStatus updates the status of an employee (used for non-activating or reactivating employee status)
func (s *Storage) UpdateEmployeeStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE employees SET status = $1 WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, status, id)
	return err
}

func (s *Storage) SelectEmployeeByID(ctx context.Context, id int) (*Employee, error) {
	var rec Employee
	query := `SELECT id, name, status, created_at FROM employees WHERE id = $1`
	err := s.db.GetContext(ctx, &rec, query, id)
	return &rec, err
}
