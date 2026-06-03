package repository

import (
	"context"
	"database/sql"
)

type queryRunner interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func selectRows[T any](
	ctx context.Context,
	query queryRunner,
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
		var row T

		row, err = scanRow(rows)

		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
