package repository

import (
	"context"
	"database/sql"
	"time"

	mssql "github.com/microsoft/go-mssqldb"
)

type clientRows struct {
	ClientCode string `tvp:"client_code"`
}

func (r *Repository) makeClientCodeList(clientCodes []string) (mssql.TVP, bool) {
	if len(clientCodes) == 0 {
		return mssql.TVP{}, false
	}

	clients := make([]clientRows, 0, len(clientCodes))
	for _, code := range clientCodes {
		clients = append(clients, clientRows{ClientCode: code})
	}
	return mssql.TVP{
		TypeName: "api.client_code_list",
		Value:    clients,
	}, true
}

type limitFilterSQL struct {
	countByClients    string
	countByClientsOp  string
	countAll          string
	countAllOp        string
	selectByClients   string
	selectByClientsOp string
	selectAll         string
	selectAllOp       string
}

func selectLimitsWithFilters[T any](
	r *Repository,
	ctx context.Context,
	opName string,
	date time.Time,
	limit uint32, offset uint64,
	clientCodes []string,
	includeTotalCount bool,
	q limitFilterSQL,
	scanRow func(*sql.Rows) (T, error),
) (result []T, totalCount *uint64, err error) {
	defer func() { err = r.finalizeSelectErr(opName, date, err) }()

	clients, hasClients := r.makeClientCodeList(clientCodes)

	query := queryRunner(r.Db)

	if hasClients {
		if includeTotalCount {
			start := time.Now()
			err = query.QueryRowContext(ctx, q.countByClients, date, sql.Named("codes", clients)).Scan(&totalCount)
			r.metrics.ObserveDBQuery(q.countByClientsOp, time.Since(start), err)
			if err != nil {
				return nil, nil, err
			}
			if *totalCount == 0 {
				return result, totalCount, err
			}

		}
		result, err = selectLimitRows(query, r, ctx, q.selectByClientsOp, q.selectByClients, scanRow, date, offset, limit, sql.Named("codes", clients))
		return result, totalCount, err
	}

	if includeTotalCount {
		start := time.Now()
		err = query.QueryRowContext(ctx, q.countAll, date).Scan(&totalCount)
		r.metrics.ObserveDBQuery(q.countAllOp, time.Since(start), err)
		if err != nil {
			return nil, nil, err
		}
		if *totalCount == 0 {
			return result, totalCount, err
		}
	}
	result, err = selectLimitRows(query, r, ctx, q.selectAllOp, q.selectAll, scanRow, date, offset, limit)
	return result, totalCount, err
}

type queryRunner interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func selectLimitRows[T any](
	query queryRunner,
	r *Repository,
	ctx context.Context,
	op string,
	sqlText string,
	scanRow func(*sql.Rows) (T, error),
	args ...any,
) (result []T, err error) {
	start := time.Now()
	rows, err := query.QueryContext(ctx, sqlText, args...)
	if err != nil {
		r.metrics.ObserveDBQuery(op, time.Since(start), err)
		return nil, err
	}
	defer rows.Close()
	defer func() { r.metrics.ObserveDBQuery(op, time.Since(start), err) }()

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
