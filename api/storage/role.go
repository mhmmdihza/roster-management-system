package storage

import (
	"context"
)

type Role struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func (s *Storage) SelectAllRoles(ctx context.Context) ([]Role, error) {
	var roles []Role
	query := `SELECT id, name FROM roles`
	err := s.db.SelectContext(ctx, &roles, query)
	return roles, err
}
