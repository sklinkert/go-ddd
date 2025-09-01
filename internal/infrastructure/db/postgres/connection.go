package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	db "github.com/sklinkert/go-ddd/internal/infrastructure/db/sqlc"
)

func NewConnection(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func NewQueries(conn *pgx.Conn) *db.Queries {
	return db.New(conn)
}
