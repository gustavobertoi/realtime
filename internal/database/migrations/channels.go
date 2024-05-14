package migrations

import (
	"database/sql"

	"github.com/gustavobertoi/realtime/pkg/logs"
)

type ChannelMigrations struct {
	db     *sql.DB
	logger *logs.Logger
}

func NewChannelMigrations(db *sql.DB, logger *logs.Logger) *ChannelMigrations {
	return &ChannelMigrations{
		db:     db,
		logger: logger,
	}
}

func (cm *ChannelMigrations) Up() error {
	cm.logger.Info("Creating channels table")
	_, err := cm.db.Exec(`CREATE TABLE IF NOT EXISTS channels (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL CHECK (type IN ('WS', 'SSE')),
		name TEXT NOT NULL,
		config TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}
	cm.logger.Info("Channels table created successfully")
	return nil
}

func (cm *ChannelMigrations) Down() error {
	cm.logger.Info("Dropping channels table")
	_, err := cm.db.Exec(`DROP TABLE channels`)
	if err != nil {
		return err
	}
	cm.logger.Info("Channels table dropped successfully")
	return nil
}
