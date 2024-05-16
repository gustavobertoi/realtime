package database

import (
	"context"
	"io/fs"
	"os"
	"path"

	"github.com/jackc/pgx/v5"
	"github.com/theritikchoure/logx"
)

type Migrations struct {
	conn *pgx.Conn
	path string
}

func NewMigrations(conn *pgx.Conn) (*Migrations, error) {
	homePath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	migrationsPath := path.Join(homePath, "internal/database/migrations")
	return &Migrations{
		conn: conn,
		path: migrationsPath,
	}, nil
}

func (m *Migrations) Init(ctx context.Context) error {
	logx.Log("Creating migrations table", logx.FGRED, logx.BGGREEN)
	_, err := m.conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		logx.Logf("Error creating migrations table: %s", logx.FGRED, logx.BGWHITE, err.Error())
		return err
	}
	logx.Log("Migrations table created", logx.FGRED, logx.BGGREEN)
	return nil
}

func (m *Migrations) Up(ctx context.Context) error {
	files, err := os.ReadDir(m.path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			return err
		}
		logx.Logf("Applying migration of %s", logx.FGRED, logx.BGGREEN, info.Name())
		sqlPath := path.Join(m.path, info.Name(), "up.sql")
		if err := m.applySql(ctx, sqlPath); err != nil {
			return err
		}
		if err := m.addMigration(ctx, info.Name()); err != nil {
			return err
		}
		logx.Logf("Migration %s applied", logx.FGRED, logx.BGGREEN, info.Name())
	}
	return nil
}

func (m *Migrations) Down(ctx context.Context) error {
	files, err := os.ReadDir(m.path)
	if err != nil {
		return err
	}

	dirs := make([]fs.FileInfo, 0)

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			return err
		}
		dirs = append(dirs, info)
	}

	reverseArrayOfMigrations(dirs)

	for _, dir := range dirs {
		logx.Logf("Rolling back migration %s", logx.FGRED, logx.BGGREEN, dir.Name())
		sqlPath := path.Join(m.path, dir.Name(), "down.sql")
		if err := m.applySql(ctx, sqlPath); err != nil {
			return err
		}
		if err := m.removeMigration(ctx, dir.Name()); err != nil {
			return err
		}
		logx.Logf("Migration %s destroyed", logx.FGRED, logx.BGGREEN, dir.Name())
	}

	return nil
}

func (m *Migrations) applySql(ctx context.Context, sqlPath string) error {
	content, err := os.ReadFile(sqlPath)
	if err != nil {
		return err
	}
	if _, err := m.conn.Exec(ctx, string(content)); err != nil {
		return err
	}
	return nil
}

func (m *Migrations) addMigration(ctx context.Context, name string) error {
	_, err := m.conn.Exec(ctx, `INSERT INTO migrations (name) VALUES ($1)`, name)
	if err != nil {
		return err
	}
	return nil
}

func (m *Migrations) removeMigration(ctx context.Context, name string) error {
	_, err := m.conn.Exec(ctx, `DELETE FROM migrations WHERE name = $1`, name)
	if err != nil {
		return err
	}
	return nil
}

func reverseArrayOfMigrations(arr []fs.FileInfo) {
	left, right := 0, len(arr)-1
	for left < right {
		arr[left], arr[right] = arr[right], arr[left]
		left++
		right--
	}
}
