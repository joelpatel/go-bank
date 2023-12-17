package db

import (
	"context"
	"database/sql"
)

// basic raw database operations
type Ops interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
}

// provides basic raw database operations
type Queries struct {
	db Ops
}

// generate queries methods for db or tx operations
func NewQueries(db Ops) *Queries {
	return &Queries{db: db}
}
