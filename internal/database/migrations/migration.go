package migrations

import (
	"database/sql"

	"github.com/gustavobertoi/realtime/pkg/logs"
)

func MigrationUp(db *sql.DB) error {
	logger := logs.NewLogger("migration:up")

	cm := NewChannelMigrations(db, logger)
	if err := cm.Up(); err != nil {
		return err
	}

	return nil
}

func MigrationDown(db *sql.DB) error {
	logger := logs.NewLogger("migration:down")

	cm := NewChannelMigrations(db, logger)
	if err := cm.Down(); err != nil {
		return err
	}

	return nil
}
