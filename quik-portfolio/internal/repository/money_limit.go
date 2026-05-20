package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
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

func (r *Repository) SelectMoneyLimits(ctx context.Context, date time.Time) (result []quik.MoneyLimit, err error) {
	defer func() { err = r.finalizeSelectErr("SelectMoneyLimits", date, err) }()

	rows, err := r.Db.QueryContext(ctx, selectMoneyLimitsByDate, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		row, err := r.scanMoneyLimit(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	if err = rows.Err(); err != nil {
		return nil, err
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
ORDER BY li.load_date, 
	li.client_code, 
	li.ccy,
	li.position_code,
    li.settle_code,
    li.firm_code DESC 
OFFSET @p2 ROWS 
FETCH NEXT @p3 ROWS ONLY`

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
ORDER BY li.load_date, 
	li.client_code, 
	li.ccy,
	li.position_code,
    li.settle_code,
    li.firm_code DESC 
OFFSET @p2 ROWS 
FETCH NEXT @p3 ROWS ONLY`

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

func (r *Repository) SelectMoneyLimitsWithFilters(ctx context.Context, date time.Time, limit, offset int, clientCodes []string) (result []quik.MoneyLimit, totalCount int, err error) {

	defer func() { err = r.finalizeSelectErr("SelectMoneyLimitsWithFilters", date, err) }()

	clients, hasClients := r.makeClientCodeList(clientCodes)

	tx, err := r.Db.BeginTx(ctx, nil)

	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	var rows *sql.Rows
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()

	if hasClients {
		err = tx.QueryRowContext(ctx, countMoneyLimitsByClients, date, sql.Named("codes", clients)).Scan(&totalCount)
		if err != nil {
			return nil, 0, err
		}

		rows, err = tx.QueryContext(ctx, selectMoneyLimitsByClients, date, offset, limit, sql.Named("codes", clients))
		if err != nil {
			return nil, 0, err
		}

	} else {
		err = tx.QueryRowContext(ctx, countMoneyLimitsAllClients, date).Scan(&totalCount)
		if err != nil {
			return nil, 0, err
		}
		rows, err = tx.QueryContext(ctx, selectMoneyLimitsAllClients, date, offset, limit)
		if err != nil {
			return nil, 0, err
		}

	}

	for rows.Next() {
		var row quik.MoneyLimit
		row, err = r.scanMoneyLimit(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, row)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	return result, totalCount, nil

}
