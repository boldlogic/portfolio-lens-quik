package repository

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

type rawInstrument struct {
	InstrumentClass        string
	SecCode                string
	TradePointId           uint8
	ISIN                   sql.NullString
	RegistrationNumber     sql.NullString
	FullName               sql.NullString
	ShortName              sql.NullString
	FaceValue              *decimal.Decimal
	MaturityDate           sql.NullTime
	CouponDuration         sql.NullInt32
	ClassCode              string
	BoardId                uint8
	Currency               sql.NullString
	CurrencyNumeric        sql.NullInt16
	BaseCurrency           sql.NullString
	BaseCurrencyNumeric    sql.NullInt16
	QuoteCurrency          sql.NullString
	QuoteCurrencyNumeric   sql.NullInt16
	CounterCurrency        sql.NullString
	CounterCurrencyNumeric sql.NullInt16
	InstrumentId           sql.NullInt64
}

const (
	SelectNewInstrumentsBatch = `
with stock_boards as (
select b.board_id,b.code, b.trade_point_id, t.type_id from ref.boards b 
join  ref.instrument_type_boards tb on b.board_id=tb.board_id
join ref.instrument_types t on t.type_id=tb.type_id where t.title=@p1
),
stocks as (

SELECT TOP (@p2)
			q.instrument_class,
			q.sec_code,
			b.trade_point_id,
			q.isin,
			q.registration_number,
			q.full_name,
			q.short_name,
			currency=norm_currency.currency,
			base_currency=norm_base_currency.base_currency,
			quote_currency=norm_quote_currency.quote_currency,
			counter_currency=norm_counter_currency.counter_currency,
			q.face_value,
			q.maturity_date,
			q.coupon_duration,
			class_code=b.code,
			b.board_id
		FROM quik.current_quotes q 
		inner join stock_boards b on b.code=q.class_code
		CROSS APPLY (
				SELECT currency = CASE
					WHEN q.currency IN ('SUR', 'RUR') THEN 'RUB'
					WHEN q.currency = 'USDX' THEN 'USD'
					ELSE q.currency
				END
			) norm_currency
		CROSS APPLY (
				SELECT base_currency = CASE
					WHEN q.base_currency IN ('SUR', 'RUR') THEN 'RUB'
					WHEN q.base_currency = 'USDX' THEN 'USD'
					ELSE q.base_currency
				END
			) norm_base_currency
		CROSS APPLY (
				SELECT quote_currency = CASE
					WHEN q.quote_currency IN ('SUR', 'RUR') THEN 'RUB'
					WHEN q.quote_currency = 'USDX' THEN 'USD'
					ELSE q.quote_currency
				END
			) norm_quote_currency
		CROSS APPLY (
				SELECT counter_currency = CASE
					WHEN q.counter_currency IN ('SUR', 'RUR') THEN 'RUB'
					WHEN q.counter_currency = 'USDX' THEN 'USD'
					ELSE q.counter_currency
				END
			) norm_counter_currency
		where q.instrument_id is null
		--and q.rw>@p2
		order by row_number() over (partition by q.sec_code,
			b.trade_point_id order by q.rw) 
			), cur_keys as (
			select distinct (currency) from stocks where currency is not null
			union
			select distinct (base_currency) from stocks where base_currency is not null
			union
			select distinct (quote_currency) from stocks where quote_currency is not null
			union
			select distinct (counter_currency) from stocks where counter_currency is not null
			), cur as (
			select k.currency,iso=coalesce(c_ec.iso_code,c_iso.iso_code)
			from cur_keys k
			LEFT JOIN ref.external_codes ec_mv
			ON ec_mv.ext_system_id = (select
				ext_system_id
			from
				ref.external_systems
			where
				ext_system = 'QUIK')
			AND ec_mv.ext_code_type_id = 1
				AND ec_mv.ext_code = k.currency
		LEFT JOIN ref.currencies c_ec  ON c_ec.iso_code  = ec_mv.internal_id
			LEFT JOIN ref.currencies c_iso ON c_iso.iso_char_code = k.currency
			)
			select 
			s.instrument_class,	
			s.sec_code,	
			s.trade_point_id,	
			s.isin,	
			s.registration_number,	
			s.full_name,
			s.short_name,	
			s.currency,	
			currency_iso=c1.iso,
			s.base_currency,
			base_currency_iso=c2.iso,
			s.quote_currency,
			quote_currency_iso=c3.iso,
			s.counter_currency,	
			counter_currency_iso=c4.iso,
			s.face_value,	
			s.maturity_date,	
			s.coupon_duration,	
			s.class_code,	
			s.board_id

from stocks s
left join cur c1 on s.currency=c1.currency
left join cur c2 on s.base_currency=c2.currency
left join cur c3 on s.quote_currency=c3.currency
left join cur c4 on s.counter_currency=c4.currency`
)

func scanInstrument(row *sql.Rows) (rawInstrument, error) {
	res := rawInstrument{}
	err := row.Scan(
		&res.InstrumentClass,
		&res.SecCode,
		&res.TradePointId,
		&res.ISIN,
		&res.RegistrationNumber,
		&res.FullName,
		&res.ShortName,
		&res.Currency,
		&res.CurrencyNumeric,
		&res.BaseCurrency,
		&res.BaseCurrencyNumeric,
		&res.QuoteCurrency,
		&res.QuoteCurrencyNumeric,
		&res.CounterCurrency,
		&res.CounterCurrencyNumeric,
		&res.FaceValue,
		&res.MaturityDate,
		&res.CouponDuration,
		&res.ClassCode,
		&res.BoardId,
	)
	if err != nil {
		return rawInstrument{}, err
	}
	return res, nil
}

// func (r *Repository) SelectNewInstruments(ctx context.Context, batchSize int, typeTitle quik.InstrumentTypeTitle) ([]models.InstrumentWithBoard, error) {
// 	var res []models.InstrumentWithBoard

// }
