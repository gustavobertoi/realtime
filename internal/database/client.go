package database

import (
	"context"
	"os"
	"strings"

	sqlc "github.com/gustavobertoi/realtime/internal/database/gen"
	"github.com/jackc/pgx/v5"
)

type DatabaseClient struct {
	conn *pgx.Conn
}

func NewClient(ctx context.Context) (*DatabaseClient, error) {
	url := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if url == "" {
		panic("DATABASE_URL is not set")
	}
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}
	return &DatabaseClient{conn: conn}, nil
}

func (db *DatabaseClient) Queries() *sqlc.Queries {
	return sqlc.New(db.conn)
}

func (db *DatabaseClient) Close(ctx context.Context) {
	db.conn.Close(ctx)
}

func (db *DatabaseClient) GetConn() *pgx.Conn {
	return db.conn
}
