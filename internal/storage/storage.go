package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDBConn(appCtx context.Context, connStr string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(appCtx, 30*time.Second)
	defer cancel()
	conn, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}
	return conn, nil
}
