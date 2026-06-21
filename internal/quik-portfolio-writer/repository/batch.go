package repository

import (
	"context"

	"database/sql"

	"github.com/boldlogic/portfolio-lens-quik/pkg/dbrepo"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
)

type moneyLimitTVP struct {
	ClientCode   string
	CurrencyCode string
	PositionCode string
	SettleCode   string
	FirmCode     string
	Balance      decimal.Decimal
}

type securityLimitTVP struct {
	ClientCode              string
	SecCode                 string
	TradeAccount            string
	SettleCode              string
	FirmCode                string
	Balance                 decimal.Decimal
	AcquisitionCurrencyCode *string
	ISIN                    *string
}

func limitToSLTVP(limit quik.Limit) securityLimitTVP {

	return securityLimitTVP{
		ClientCode:              limit.ClientCode(),
		SecCode:                 limit.Ticker(),
		TradeAccount:            *limit.TradeAccount(),
		SettleCode:              limit.SettleCode().String(),
		FirmCode:                limit.FirmCode(),
		Balance:                 limit.Balance(),
		AcquisitionCurrencyCode: limit.AcquisitionCurrencyCode(),
		ISIN:                    limit.ISIN(),
	}
}

func limitToMLTVP(limit quik.Limit) moneyLimitTVP {
	return moneyLimitTVP{
		ClientCode:   limit.ClientCode(),
		CurrencyCode: limit.Ticker(),
		PositionCode: *limit.PositionCode(),
		SettleCode:   limit.SettleCode().String(),
		FirmCode:     limit.FirmCode(),
		Balance:      limit.Balance(),
	}
}

const (
	mergeSecurityLimits           = mergeSecurityLimitsSrcSQL + securityLimitExchangeTableSQL + mergeSecurityLimitsUpsertSQL
	mergeSecurityLimitsOTC        = mergeSecurityLimitsSrcSQL + securityLimitOtcTableSQL + mergeSecurityLimitsUpsertSQL
	securityLimitExchangeTableSQL = " quik.security_limits "
	securityLimitOtcTableSQL      = " quik.security_limits_otc "
	mergeSecurityLimitsSrcSQL     = `
	WITH src as(
		select client_code
				,sec_code
				,trade_account
				,settle_code
				,firm_code
				,balance
				,acquisition_currency_code
				,isin
		from @limits)
		merge into `
	mergeSecurityLimitsUpsertSQL = ` as tgt
		using src on 
			tgt.load_date=cast(GETDATE() as date)
			and tgt.client_code=src.client_code
			and tgt.sec_code=src.sec_code
			and tgt.trade_account=src.trade_account
			and tgt.settle_code=src.settle_code
			and tgt.firm_code=src.firm_code
		when matched 
			and (tgt.balance<>src.balance
				or tgt.acquisition_currency_code<>src.acquisition_currency_code
				or (tgt.acquisition_currency_code is null and src.acquisition_currency_code is not null)
				or (tgt.acquisition_currency_code is not null and src.acquisition_currency_code is null)
				or tgt.isin<>src.isin
				or (tgt.isin is null and src.isin is not null)
				or (tgt.isin is not null and src.isin is null))	
			then update set
				tgt.balance=src.balance,
				tgt.acquisition_currency_code=src.acquisition_currency_code,
				tgt.isin=src.isin
		WHEN NOT MATCHED BY TARGET
			then insert	(
					client_code, 
					sec_code, 
					trade_account, 
					settle_code, 
					firm_code, 
					balance,
					acquisition_currency_code,
					isin)	
				values(
					src.client_code,
					src.sec_code,
					src.trade_account,
					src.settle_code,
					src.firm_code,
					src.balance,
					src.acquisition_currency_code,
					src.isin
					);
	`

	mergeMoneyLimits = `
	WITH src as(
		select client_code,
				currency_code,
				position_code,
				settle_code,
				firm_code,
				balance
		from @limits)
	merge into quik.money_limits as tgt
	using src on 
		tgt.load_date=cast(GETDATE() as date)
		and tgt.client_code=src.client_code
		and tgt.currency_code=src.currency_code
		and tgt.position_code=src.position_code
		and tgt.settle_code=src.settle_code
		and tgt.firm_code=src.firm_code
	when matched 
		and tgt.balance<>src.balance	
		then update set
			tgt.balance=src.balance
	WHEN NOT MATCHED BY TARGET
		then insert	(
				client_code, 
				currency_code, 
				position_code, 
				settle_code, 
				firm_code, 
				balance)	
			values(
				src.client_code,
				src.currency_code,
				src.position_code,
				src.settle_code,
				src.firm_code,
				src.balance
				);
	`
)

func (r *Repository) HandleRequest(ctx context.Context, limits []quik.Limit) error {

	ml := make([]quik.Limit, 0, len(limits))
	sl := make([]quik.Limit, 0, len(limits))
	slo := make([]quik.Limit, 0, len(limits))
	for _, limit := range limits {
		switch limit.Type() {
		case quik.LimitTypeMoney:
			ml = append(ml, limit)
		case quik.LimitTypeSecurities:
			sl = append(sl, limit)
		case quik.LimitTypeSecuritiesOtc:
			slo = append(slo, limit)
		}
	}
	mlT, okM := dbrepo.MakeTVP(ml, limitToMLTVP, "app.money_limits")
	slT, okS := dbrepo.MakeTVP(sl, limitToSLTVP, "app.security_limits")
	slOT, okO := dbrepo.MakeTVP(slo, limitToSLTVP, "app.security_limits")
	g, gCTX := errgroup.WithContext(ctx)
	if okM {
		g.Go(func() error {
			_, err := r.Db.ExecContext(gCTX, mergeMoneyLimits, sql.Named("limits", mlT))
			return err
		})
	}
	if okS {
		g.Go(func() error {
			_, err := r.Db.ExecContext(gCTX, mergeSecurityLimits, sql.Named("limits", slT))
			return err
		})
	}
	if okO {
		g.Go(func() error {
			_, err := r.Db.ExecContext(gCTX, mergeSecurityLimitsOTC, sql.Named("limits", slOT))
			return err
		})
	}
	err := g.Wait()
	if err != nil {
		r.Logger.Error(err.Error())
		return err
	}
	return nil

}
