package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
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
		TypeName: "app.client_code_list",
		Value:    clients,
	}, true
}

type limitFilterSQL struct {
	countByClients  string
	countAll        string
	selectByClients string
	selectAll       string
}

type limitListQuery struct {
	date              time.Time
	limit             uint32
	offset            uint64
	clientCodes       []string
	includeTotalCount bool
}

type portfolioQuery struct {
	date      time.Time
	targetCcy string
}

func limitFilterSQLByType(limitType quik.LimitType) (limitFilterSQL, error) {
	switch limitType {
	case quik.LimitTypeMoney:
		return limitFilterSQL{
			countByClients:  countMoneyLimitsByClients,
			countAll:        countMoneyLimitsAllClients,
			selectByClients: selectMoneyLimitsByClients,
			selectAll:       selectMoneyLimitsAllClients,
		}, nil
	case quik.LimitTypeSecurities:
		return limitFilterSQL{
			countByClients:  countSecurityLimitsByClients,
			countAll:        countSecurityLimitsAllClients,
			selectByClients: selectSecurityLimitsByClients,
			selectAll:       selectSecurityLimitsAllClients,
		}, nil
	case quik.LimitTypeSecuritiesOtc:
		return limitFilterSQL{
			countByClients:  countSecurityLimitsOtcByClients,
			countAll:        countSecurityLimitsOtcAllClients,
			selectByClients: selectSecurityLimitsOtcByClients,
			selectAll:       selectSecurityLimitsOtcAllClients,
		}, nil
	default:
		return limitFilterSQL{}, fmt.Errorf("неподдерживаемый тип лимита: %s", limitType)
	}
}

func selectLimitRows[T any](
	r *Repository,
	ctx context.Context,
	opName string,
	query limitListQuery,
	limitType quik.LimitType,
	scanRow func(*sql.Rows) (T, error),
) (result []T, totalCount *uint64, err error) {
	defer func() { err = r.finalizeSelectErr(opName, query.date, err) }()

	q, err := limitFilterSQLByType(limitType)
	if err != nil {
		return nil, nil, err
	}

	clients, hasClients := r.makeClientCodeList(query.clientCodes)

	db := queryRunner(r.Db)

	if hasClients {
		if query.includeTotalCount {
			err = db.QueryRowContext(ctx, q.countByClients, query.date, sql.Named("codes", clients)).Scan(&totalCount)
			if err != nil {
				return nil, nil, err
			}
			if *totalCount == 0 {
				return result, totalCount, err
			}
		}
		result, err = selectRows(
			ctx,
			db,
			q.selectByClients,
			scanRow,
			query.date,
			query.offset,
			query.limit,
			sql.Named("codes", clients))
		return result, totalCount, err
	}

	if query.includeTotalCount {
		err = db.QueryRowContext(ctx, q.countAll, query.date).Scan(&totalCount)
		if err != nil {
			return nil, nil, err
		}
		if *totalCount == 0 {
			return result, totalCount, err
		}
	}
	result, err = selectRows(ctx, db, q.selectAll, scanRow, query.date, query.offset, query.limit)
	return result, totalCount, err
}

func selectPortfolioRows[T any](
	r *Repository,
	ctx context.Context,
	opName string,
	sqlText string,
	scanRow func(*sql.Rows) (T, error),
	query portfolioQuery,
) (result []T, err error) {
	start := time.Now()
	defer func() { r.metrics.ObserveRepository(opName, time.Since(start), err) }()

	return selectRows(ctx, r.Db, sqlText, scanRow, query.date, query.targetCcy)
}
