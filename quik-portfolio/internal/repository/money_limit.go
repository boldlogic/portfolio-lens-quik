package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	mssql "github.com/microsoft/go-mssqldb"
	"go.uber.org/zap"
)

const (
	selectMoneyLimitsByDate = `
		WITH cte AS (
			SELECT
				li.load_date,
				li.source_date,
				li.client_code,
				li.ccy,
				li.position_code,
				li.firm_code,
				li.settle_code,
				li.firm_name,
				li.balance,
				settle_max = MAX(li.settle_code) OVER (
					PARTITION BY li.load_date, li.client_code, li.ccy, li.position_code, li.firm_code
				)
			FROM quik.money_limits li
			WHERE li.load_date = cast(@p1 as date)
		)
		SELECT
			c.load_date,
			c.source_date,
			c.client_code,
			c.ccy,
			c.position_code,
			c.settle_code,
			c.firm_code,
			c.firm_name,
			c.balance
		FROM cte c
		WHERE 1=1
		ORDER BY c.load_date, c.client_code, c.ccy, c.position_code, c.firm_code;
`
)

func (r *Repository) SelectMoneyLimits(ctx context.Context, date time.Time) ([]quik.MoneyLimit, error) {
	var result []quik.MoneyLimit

	rows, err := r.Db.QueryContext(ctx, selectMoneyLimitsByDate, date)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка выполнения запроса на получение позиций по деньгам", zap.Time("load_date", date), zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	defer rows.Close()

	for rows.Next() {
		row, err := r.scanMoneyLimit(rows)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}
			r.Logger.Error("ошибка чтения позиций по деньгам", zap.Time("load_date", date), zap.Error(err))
			return nil, models.ErrRetrievingData
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		r.Logger.Error("ошибка чтения позиций по деньгам", zap.Time("load_date", date), zap.Error(rows.Err()))
		return nil, models.ErrRetrievingData
	}

	r.Logger.Debug("результаты получения позиций по деньгам", zap.Time("load_date", date), zap.Int("count", len(result)))

	if len(result) == 0 {
		r.Logger.Warn("позиции по деньгам не найдены", zap.Time("load_date", date))
	}
	return result, nil
}

const selectMoneyLimitsByClients = `SELECT
    li.load_date,
    li.source_date,
    li.client_code,
    li.ccy,
    li.position_code,
    li.settle_code,
    li.firm_code,
    li.firm_name,
    li.balance
FROM quik.money_limits li
join @codes c on c.client_code = li.client_code
WHERE 1 = 1
    and li.load_date = case 
        when isnull(@p1, '') = '' then cast(getdate() as date) 
        else cast(@p1 as date) 
        end
ORDER BY 1, 2, 3,4 DESC 
OFFSET @p2 ROWS FETCH NEXT @p3 ROWS ONLY`

const selectMoneyLimitsAllClients = `SELECT
    li.load_date,
    li.source_date,
    li.client_code,
    li.ccy,
    li.position_code,
    li.settle_code,
    li.firm_code,
    li.firm_name,
    li.balance
FROM quik.money_limits li
WHERE 1 = 1
    and li.load_date = case 
        when isnull(@p1, '') = '' then cast(getdate() as date) 
        else cast(@p1 as date) 
        end
ORDER BY 1, 2, 3,4 DESC 
OFFSET @p2 ROWS FETCH NEXT @p3 ROWS ONLY`

const countMoneyLimitsByClients = `SELECT COUNT(*)
FROM quik.money_limits li
join @codes c on c.client_code = li.client_code
WHERE li.load_date = case 
    when isnull(@p1, '') = '' then cast(getdate() as date) 
    else cast(@p1 as date) 
    end`

const countMoneyLimitsAllClients = `SELECT COUNT(*)
FROM quik.money_limits li
WHERE li.load_date = case 
    when isnull(@p1, '') = '' then cast(getdate() as date) 
    else cast(@p1 as date) 
    end`

func (r *Repository) scanMoneyLimit(row *sql.Rows) (quik.MoneyLimit, error) {
	out := quik.MoneyLimit{}
	err := row.Scan(
		&out.LoadDate,
		&out.SourceDate,
		&out.ClientCode,
		&out.Currency,
		&out.PositionCode,
		&out.SettleCode,
		&out.FirmCode,
		&out.FirmName,
		&out.Balance,
	)
	if err != nil {
		return quik.MoneyLimit{}, err
	}
	return out, nil
}

func (r *Repository) SelectMoneyLimitsWithFilters(ctx context.Context, date time.Time, limit, offset int, clientCodes []string) ([]quik.MoneyLimit, int, error) {
	var result []quik.MoneyLimit

	clients, ok := r.makeClientCodeList(clientCodes)
	totalCount, err := r.countMoneyLimitsWithFilters(ctx, date, clients, ok)
	if err != nil {
		return nil, 0, err
	}

	var rows *sql.Rows
	if ok {
		rows, err = r.Db.QueryContext(ctx, selectMoneyLimitsByClients, date, offset, limit, sql.Named("codes", clients))
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, 0, err
			}
			r.Logger.Error("ошибка выполнения запроса на получение позиций по деньгам", zap.Time("load_date", date), zap.Error(err))
			return nil, 0, models.ErrRetrievingData
		}
		defer rows.Close()

	} else {
		rows, err = r.Db.QueryContext(ctx, selectMoneyLimitsAllClients, date, offset, limit)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, 0, err
			}
			r.Logger.Error("ошибка выполнения запроса на получение позиций по деньгам", zap.Time("load_date", date), zap.Error(err))
			return nil, 0, models.ErrRetrievingData
		}
		defer rows.Close()
	}

	for rows.Next() {
		row, err := r.scanMoneyLimit(rows)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, 0, err
			}
			r.Logger.Error("ошибка чтения позиций по деньгам", zap.Time("load_date", date), zap.Error(err))
			return nil, 0, models.ErrRetrievingData
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		r.Logger.Error("ошибка чтения позиций по деньгам", zap.Time("load_date", date), zap.Error(rows.Err()))
		return nil, 0, models.ErrRetrievingData
	}
	return result, totalCount, nil

}

func (r *Repository) countMoneyLimitsWithFilters(ctx context.Context, date time.Time, clients mssql.TVP, hasClients bool) (int, error) {
	var totalCount int
	var err error
	if hasClients {
		err = r.Db.QueryRowContext(ctx, countMoneyLimitsByClients, date, sql.Named("codes", clients)).Scan(&totalCount)
	} else {
		err = r.Db.QueryRowContext(ctx, countMoneyLimitsAllClients, date).Scan(&totalCount)
	}
	if err != nil {
		if shutdown.IsExceeded(err) {
			return 0, err
		}
		r.Logger.Error("ошибка выполнения запроса на подсчет позиций по деньгам", zap.Time("load_date", date), zap.Error(err))
		return 0, models.ErrRetrievingData
	}
	return totalCount, nil
}

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
