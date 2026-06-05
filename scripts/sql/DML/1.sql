declare @p1 date='2026-06-04';
declare @p2 varchar(3)='RUB';
WITH cte AS (
  SELECT
            li.load_date,
            li.source_date,
            li.client_code,
            li.ticker,
			li.sec_name,
            li.settle_code,
			li.firm_code,
            li.firm_name,
            li.balance,
			--acquisition_currency_code = case when UPPER(TRIM(li.acquisition_ccy)) in ('SUR', 'RUR') THEN 'RUB' ELSE UPPER(TRIM(li.acquisition_ccy)) END,
            li.acquisition_ccy,
			li.isin,
            settle_max = MAX(li.settle_code) OVER (
                PARTITION BY li.load_date, li.client_code, li.ticker, li.trade_account, li.firm_code
            )
        FROM quik.security_limits li
        WHERE li.load_date = cast(@p1 as date)
		)
		select c.load_date,
		c.source_date,
		q.quote_date,
		c.client_code,
		c.firm_code,
		c.firm_name,
		c.ticker,
		c.sec_name,
		c.balance,
		c.ticker,
		q.price_in_ccy,
		q.accrued_int,
		q.short_name,
        c.acquisition_ccy,
		q.mv_currency,
		norm.norm_mv,
		norm.norm_acc,
		q.accrued_currency
		,f_mv.rate
		,f_acc.rate
		--acquisition_currency_code=COALESCE(c_acq_ec.iso_char_code,   c_acq_iso.iso_char_code,c.acquisition_currency_code)
		,ft.rate
		,market_value=(		(isnull(q.price_in_ccy,0)*c.balance*f_mv.rate+isnull(q.accrued_int,0)*c.balance*f_acc.rate)/ft.rate		)
		from cte c
		outer apply (select price_in_ccy,accrued_int,short_name,quote_date,mv_currency,accrued_currency from dbo.fnGetQuoteByAcquisitionCurrency(c.ticker,c.acquisition_ccy)) q
		CROSS APPLY (VALUES (
    CASE WHEN UPPER(TRIM(ISNULL(q.mv_currency,      ''))) IN ('SUR','RUR') THEN 'RUB' ELSE UPPER(TRIM(ISNULL(q.mv_currency,      ''))) END,
    CASE WHEN UPPER(TRIM(ISNULL(q.accrued_currency, ''))) IN ('SUR','RUR') THEN 'RUB' ELSE UPPER(TRIM(ISNULL(q.accrued_currency, ''))) END,
    CASE WHEN UPPER(TRIM(ISNULL(@p2, 'RUB')))             IN ('SUR','RUR') THEN 'RUB' ELSE UPPER(TRIM(ISNULL(@p2, 'RUB')))             END
)) norm(norm_mv, norm_acc,norm_tgt)
		LEFT JOIN dbo.external_codes ec_mv
			ON ec_mv.ext_system_id = (select
				ext_system_id
			from
				dbo.external_systems
			where
				ext_system = 'QUIK')
			AND ec_mv.ext_code_type_id = 1
			AND ec_mv.ext_code = norm.norm_mv
		LEFT JOIN currencies c_mv_ec  ON c_mv_ec.iso_code  = ec_mv.internal_id
		LEFT JOIN currencies c_mv_iso ON c_mv_iso.iso_char_code = norm.norm_mv

		LEFT JOIN dbo.external_codes ec_acc
			ON ec_acc.ext_system_id = (select
				ext_system_id
			from
				dbo.external_systems
			where
				ext_system = 'QUIK')
			AND ec_acc.ext_code_type_id = 1
			AND ec_acc.ext_code = norm.norm_acc
		LEFT JOIN currencies c_acc_ec  ON c_acc_ec.iso_code  = ec_acc.internal_id
		LEFT JOIN currencies c_acc_iso ON c_acc_iso.iso_char_code = norm.norm_acc

		CROSS APPLY dbo.fnFxRateToRub(ISNULL(COALESCE(c_mv_ec.iso_char_code,c_mv_iso.iso_char_code), ''), c.load_date) f_mv
				CROSS APPLY dbo.fnFxRateToRub(ISNULL(COALESCE(c_acc_ec.iso_char_code,c_acc_iso.iso_char_code), ''), c.load_date) f_acc
								CROSS APPLY dbo.fnFxRateToRub(@p2,            c.load_date) ft

		where c.settle_code=c.settle_max






