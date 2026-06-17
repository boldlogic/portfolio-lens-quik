package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/internal/models"
	erss "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	mssql "github.com/microsoft/go-mssqldb"
)

func (r *Repository) InsertLimit(ctx context.Context, l quik.Limit) (quik.Limit, error) {
	var err error
	var res quik.Limit

	switch l.Type {
	case quik.LimitTypeMoney:
		var m models.MoneyLimitRow

		row := r.Db.QueryRowContext(ctx, insertMoneyLimit,
			l.ClientCode,
			l.CurrencyCode,
			l.PositionCode,
			l.SettleCode,
			l.FirmCode,
			l.Balance,
		)
		m, err = models.ScanMoneyLimitRow(row)

		res = m.ToLimit()
	case quik.LimitTypeSecurities:
		var s models.SecurityLimitRow
		row := r.Db.QueryRowContext(ctx, insertSecurityLimit,
			l.ClientCode,
			l.SecCode,
			l.TradeAccount,
			l.SettleCode,
			l.FirmCode,
			l.Balance,
			l.AcquisitionCurrencyCode,
			l.ISIN)
		s, err = models.ScanSecurityLimitRow(row)

		res = s.ToLimit()
	case quik.LimitTypeSecuritiesOtc:
		var s models.SecurityLimitRow
		row := r.Db.QueryRowContext(ctx, insertSecurityLimitOtc,
			l.ClientCode,
			l.SecCode,
			l.TradeAccount,
			l.SettleCode,
			l.FirmCode,
			l.Balance,
			l.AcquisitionCurrencyCode,
			l.ISIN)
		s, err = models.ScanSecurityLimitRow(row)

		res = s.ToLimit()
	default:
		return quik.Limit{}, fmt.Errorf("неподдерживаемый тип лимита")
	}
	var msErr mssql.Error
	switch {
	case errors.As(err, &msErr) && (msErr.Number == 2627 || msErr.Number == 2601):
		return quik.Limit{}, erss.ErrConflict
	case shutdown.IsExceeded(err):
		return quik.Limit{}, err
	case err == nil:
		return res, nil
	default:
		return quik.Limit{}, erss.ErrSavingData
	}

}

const (
	securityLimitExchangeTableSQL = "quik.security_limits"
	securityLimitOtcTableSQL      = "quik.security_limits_otc"

	insertSecurityLimitSrcSQL = `
		WITH src AS
		(
			SELECT  client_code = @p1
				,sec_code = @p2
				,trade_account = @p3
				,settle_code = @p4
				,firm_code = @p5
				,balance = @p6
				,acquisition_currency_code = @p7
				,isin = @p8
		)
		INSERT INTO `

	insertSecurityLimitTailSQL = `
		(client_code, sec_code, trade_account, settle_code, firm_code, firm_name, balance, acquisition_currency_code, isin)
		OUTPUT inserted.load_date, inserted.source_date, inserted.client_code, inserted.sec_code, inserted.trade_account, inserted.settle_code, inserted.firm_code, inserted.firm_name, inserted.balance, inserted.acquisition_currency_code, inserted.isin,inserted.sec_name
		SELECT src.client_code
			,src.sec_code
			,src.trade_account
			,src.settle_code
			,f.code
			,f.name
			,src.balance
			,src.acquisition_currency_code
			,src.isin
		FROM src
		join ref.firms f on f.code = src.firm_code
	`

	insertSecurityLimit    = insertSecurityLimitSrcSQL + securityLimitExchangeTableSQL + insertSecurityLimitTailSQL
	insertSecurityLimitOtc = insertSecurityLimitSrcSQL + securityLimitOtcTableSQL + insertSecurityLimitTailSQL
	insertMoneyLimit       = `
		WITH src AS
		(
			SELECT  client_code = @p1
				,currency_code = @p2
				,position_code = @p3
				,settle_code = @p4
				,firm_code = @p5
				,balance = @p6
		)
		INSERT INTO quik.money_limits (client_code, currency_code, position_code, settle_code, firm_code, firm_name, balance)
		OUTPUT inserted.load_date, inserted.source_date, inserted.client_code, inserted.currency_code, inserted.position_code, inserted.settle_code, inserted.firm_code, inserted.firm_name, inserted.balance
		SELECT  src.client_code
			,src.currency_code
			,src.position_code
			,src.settle_code
			,f.code
			,f.name
			,src.balance
		FROM src
		join ref.firms f on f.code = src.firm_code	
		`
)
