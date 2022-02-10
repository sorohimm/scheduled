package interfaces

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type IDBHandler interface {
	GetPool() *pgxpool.Pool
	AcquireConn(ctx context.Context) (*pgxpool.Conn, error)
}
