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
	countByClients  string
	countAll        string
	selectByClients string
	selectAll       string
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
			err = query.QueryRowContext(ctx, q.countByClients, date, sql.Named("codes", clients)).Scan(&totalCount)
			if err != nil {
				return nil, nil, err
			}
			if *totalCount == 0 {
				return result, totalCount, err
			}
		}
		result, err = selectRows(
			ctx,
			query,
			q.selectByClients,
			scanRow,
			date,
			offset,
			limit,
			sql.Named("codes", clients))
		return result, totalCount, err
	}

	if includeTotalCount {
		err = query.QueryRowContext(ctx, q.countAll, date).Scan(&totalCount)
		if err != nil {
			return nil, nil, err
		}
		if *totalCount == 0 {
			return result, totalCount, err
		}
	}
	result, err = selectRows(ctx, query, q.selectAll, scanRow, date, offset, limit)
	return result, totalCount, err
}
