package database

import (
	"context"
	"database/sql"
)

type DBTx interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Close() error
}


func NewRead() DBTx {
	return &sql.DB{}
}

func NewWrite() DBTx {
	return &sql.DB{}
}
