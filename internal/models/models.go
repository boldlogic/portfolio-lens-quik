package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/dbrepo"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

type SecurityLimitRow struct {
	LoadDate                time.Time
	SourceDate              time.Time
	ClientCode              string
	SecCode                 string
	TradeAccount            string
	SettleCode              string
	FirmCode                string
	FirmName                sql.NullString
	Balance                 *decimal.Decimal
	AcquisitionCurrencyCode string
	ISIN                    sql.NullString
	ShortName               sql.NullString
}

type MoneyLimitRow struct {
	LoadDate     time.Time
	SourceDate   time.Time
	ClientCode   string
	CurrencyCode string
	PositionCode string
	SettleCode   string
	FirmCode     string
	FirmName     sql.NullString
	Balance      *decimal.Decimal
}

func (m MoneyLimitRow) ToLimit() quik.Limit {
	//if m.CurrencyCode
	ml := quik.Limit{
		LoadDate:     m.LoadDate,
		SourceDate:   m.SourceDate,
		ClientCode:   m.ClientCode,
		CurrencyCode: &m.CurrencyCode,
		PositionCode: &m.PositionCode,
		FirmCode:     m.FirmCode,
		FirmName:     dbrepo.StringFromNull(m.FirmName),
		Balance:      dbrepo.DecimalFromPtr(m.Balance),
		SettleCode:   quik.SettleCode(m.SettleCode),
	}
	return ml
}

func (s SecurityLimitRow) ToLimit() quik.Limit {
	var isin *string
	if s.ISIN.Valid {
		isin = &s.ISIN.String
	}
	var shortName *string
	if s.ShortName.Valid {
		isin = &s.ShortName.String
	}
	return quik.Limit{
		LoadDate:                s.LoadDate,
		SourceDate:              s.SourceDate,
		ClientCode:              s.ClientCode,
		SecCode:                 &s.SecCode,
		TradeAccount:            &s.TradeAccount,
		FirmCode:                s.FirmCode,
		FirmName:                dbrepo.StringFromNull(s.FirmName),
		Balance:                 dbrepo.DecimalFromPtr(s.Balance),
		AcquisitionCurrencyCode: &s.AcquisitionCurrencyCode,
		ISIN:                    isin,
		ShortName:               shortName,
		SettleCode:              quik.SettleCode(s.SettleCode),
	}
}

func (m MoneyLimitRow) ToQuik() quik.MoneyLimit {
	ml := quik.MoneyLimit{
		LoadDate:     m.LoadDate,
		SourceDate:   m.SourceDate,
		ClientCode:   m.ClientCode,
		CurrencyCode: m.CurrencyCode,
		PositionCode: m.PositionCode,
		FirmCode:     m.FirmCode,
		FirmName:     dbrepo.StringFromNull(m.FirmName),
		Balance:      dbrepo.DecimalFromPtr(m.Balance),
		SettleCode:   quik.SettleCode(m.SettleCode),
	}
	return ml
}

func (s SecurityLimitRow) ToQuik() quik.SecurityLimit {
	return quik.SecurityLimit{
		LoadDate:                s.LoadDate,
		SourceDate:              s.SourceDate,
		ClientCode:              s.ClientCode,
		SecCode:                 s.SecCode,
		TradeAccount:            s.TradeAccount,
		FirmCode:                s.FirmCode,
		FirmName:                dbrepo.StringFromNull(s.FirmName),
		Balance:                 dbrepo.DecimalFromPtr(s.Balance),
		AcquisitionCurrencyCode: s.AcquisitionCurrencyCode,
		ISIN:                    dbrepo.StringFromNull(s.ISIN),
		ShortName:               dbrepo.StringFromNull(s.ShortName),
		SettleCode:              quik.SettleCode(s.SettleCode),
	}
}

var (
	ErrScan = errors.New("ошибка чтения строки")
)

func ScanMoneyLimitRow(row *sql.Row) (MoneyLimitRow, error) {
	out := MoneyLimitRow{}
	err := row.Scan(
		&out.LoadDate,
		&out.SourceDate,
		&out.ClientCode,
		&out.CurrencyCode,
		&out.PositionCode,
		&out.SettleCode,
		&out.FirmCode,
		&out.FirmName,
		&out.Balance,
	)
	if err != nil {
		return MoneyLimitRow{}, fmt.Errorf("%w: %w", ErrScan, err)
	}
	return out, nil
}

func ScanSecurityLimitRow(row *sql.Row) (SecurityLimitRow, error) {
	out := SecurityLimitRow{}
	err := row.Scan(
		&out.LoadDate,
		&out.SourceDate,
		&out.ClientCode,
		&out.SecCode,
		&out.SettleCode,
		&out.TradeAccount,
		&out.FirmCode,
		&out.FirmName,
		&out.Balance,
		&out.AcquisitionCurrencyCode,
		&out.ISIN,
		&out.ShortName,
	)
	if err != nil {
		return SecurityLimitRow{}, fmt.Errorf("%w: %w", ErrScan, err)
	}

	return out, nil
}
