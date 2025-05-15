package storage

import (
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
