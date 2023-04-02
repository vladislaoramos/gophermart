package app

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

const (
	driver       = "postgres"
	migrationDir = "migrations"
)

func applyMigration(dbURL string) error {
	db, err := goose.OpenDBWithDriver(driver, dbURL)
	if err != nil {
		return fmt.Errorf("error open db: %w", err)
	}

	defer db.Close()

	if err := goose.Up(db, migrationDir); err != nil {
		return fmt.Errorf("error up migration: %w", err)
	}

	return nil
}
