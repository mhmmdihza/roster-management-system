package storage

import (
	"log"
	"os"
	"testing"
)

var _testStorage *Storage

func TestMain(m *testing.M) {
	var err error
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatalf("DB Name can't be empty")
	}
	_testStorage, err = NewStorage(
		WithUser(user),
		WithPassword(password),
		WithDbname(dbName))
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	code := m.Run()
	os.Exit(code)
}

func _withTestDatabase(t *testing.T, testFunc func(storage *Storage)) {
	err := _testStorage.RunMigrations("./migrations")
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	t.Cleanup(func() {
		err := _testStorage.CleanMigrations("./migrations", 0)
		if err != nil {
			t.Errorf("failed to clean migrations: %v", err)
		}
	})

	testFunc(_testStorage)
}

func TestRunMigrations(t *testing.T) {
	_withTestDatabase(t, func(storage *Storage) {
	})
}
