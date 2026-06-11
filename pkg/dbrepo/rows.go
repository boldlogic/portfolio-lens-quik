package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type QueryRunner interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type ExecRunner interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

var ErrScan = errors.New("ошибка чтения строки")
var ErrRows = errors.New("ошибка обхода строк")

func ScanError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %w", ErrScan, err)
}

func SelectRows[T any](
	ctx context.Context,
	query QueryRunner,
	sqlText string,
	scanRow func(*sql.Rows) (T, error),
	args ...any,
) (result []T, err error) {
	rows, err := query.QueryContext(ctx, sqlText, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		row, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRows, err)
	}
	return result, nil
}
