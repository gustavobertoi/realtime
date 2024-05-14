package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

const (
	TursoDriver = "libsql"
)

var (
	ErrDatabaseURLNotSet = errors.New("database url not set")
)

func NewTursoDB() (*sql.DB, error) {
	url := fmt.Sprintf("%s?authToken=%s", strings.TrimSpace(os.Getenv("TURSO_DATABASE_URL")), strings.TrimSpace(os.Getenv("TURSO_AUTH_TOKEN")))
	db, err := sql.Open(TursoDriver, url)
	if err != nil {
		return nil, err
	}
	return db, nil
}
