package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type Storage struct {
	db       *sqlx.DB
	user     string
	password string
	host     string
	port     int
	dbname   string
	sslmode  string
}

type StorageOption func(*Storage) error

func NewStorage(opts ...StorageOption) (*Storage, error) {
	s := &Storage{
		host:    "localhost",
		port:    5432,
		sslmode: "disable",
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		s.user, s.password, s.host, s.port, s.dbname, s.sslmode)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	s.db = db
	return s, nil
}

type pgTx struct{}

type tx struct {
	*sqlx.Tx
	committed *bool
}

// use this context to perform a transactional commit,
// allowing a rollback if the process is interrupted.
// see example on CreateNewEmployee
func (s *Storage) NewTransacton(ctx context.Context) (context.Context, error) {
	t, err := s.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, pgTx{}, &tx{
		Tx:        t,
		committed: new(bool),
	}), nil
}

func getTx(ctx context.Context) *tx {
	if t, ok := ctx.Value(pgTx{}).(*tx); ok && !*t.committed {
		return t
	}
	return nil
}

func (s *Storage) Commit(ctx context.Context) error {
	t := getTx(ctx)
	if t == nil {
		return fmt.Errorf("not a transaction")
	}
	if *t.committed {
		return nil
	}
	if err := t.Commit(); err != nil {
		return err
	}
	*t.committed = true
	return nil
}

func (s *Storage) Rollback(ctx context.Context) error {
	t := getTx(ctx)
	if t == nil {
		return fmt.Errorf("not a transaction")
	}
	if *t.committed {
		return nil
	}
	return t.Rollback()
}

func (s *Storage) RunMigrations(migrationDir string) error {
	return goose.Run("up", s.db.DB, migrationDir)
}

func (s *Storage) CleanMigrations(migrationDir string, version int64) error {
	return goose.DownTo(s.db.DB, migrationDir, version)
}

func WithSSLMode(mode string) StorageOption {
	return func(s *Storage) error {
		validModes := map[string]bool{
			"disable":     true,
			"require":     true,
			"verify-ca":   true,
			"verify-full": true,
		}
		if !validModes[mode] {
			return fmt.Errorf("invalid sslmode: %s", mode)
		}
		s.sslmode = mode
		return nil
	}
}

func WithUser(user string) StorageOption {
	return func(s *Storage) error {
		s.user = user
		return nil
	}
}

func WithPassword(password string) StorageOption {
	return func(s *Storage) error {
		s.password = password
		return nil
	}
}

func WithHost(host string) StorageOption {
	return func(s *Storage) error {
		s.host = host
		if s.host == "" {
			s.host = "localhost"
		}
		return nil
	}
}

func WithPort(port int) StorageOption {
	return func(s *Storage) error {
		s.port = port
		if s.port == 0 {
			s.port = 5432
		}
		return nil
	}
}

func WithDbname(dbname string) StorageOption {
	return func(s *Storage) error {
		s.dbname = dbname
		return nil
	}
}
