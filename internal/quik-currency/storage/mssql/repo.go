package mssql

import (
	"context"
	"database/sql"
)

type QueryRunner interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type CurrencyRepo struct {
	runner QueryRunner
}

func NewCurrencyRepo(runner QueryRunner) *CurrencyRepo {
	return &CurrencyRepo{
		runner: runner,
	}
}
