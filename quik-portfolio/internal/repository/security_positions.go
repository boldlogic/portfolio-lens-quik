package repository

import "database/sql"

const selectSecurityPositions = `
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
		cr.rate_date,
		c.client_code,
		c.firm_code,
		c.firm_name,
		c.ticker,	
		sec_name=coalesce(c.sec_name,q.short_name),
		c.balance,
		q.price,
		q.ai,
		market_value_instr=cast(q.market_value*c.balance as decimal(18,2)),
		norm_ccy.currency,
		market_value=cast (q.market_value*c.balance*cr.rate as decimal(18,2))
		from cte c
		outer apply dbo.fnGetQuoteByAcquisitionCurrency(c.ticker,c.acquisition_ccy) q
		outer apply (select ccy=coalesce(q.currency,c.acquisition_ccy)) cur
		outer apply (select currency=case when isnull(cur.ccy,'')  IN ('SUR','RUR') THEN 'RUB' ELSE cur.ccy END) norm_ccy
		LEFT JOIN dbo.external_codes ec_mv
			ON ec_mv.ext_system_id = (select
				ext_system_id
			from
				dbo.external_systems
			where
				ext_system = 'QUIK')
			AND ec_mv.ext_code_type_id = 1
			AND ec_mv.ext_code = norm_ccy.currency
		LEFT JOIN currencies c_mv_ec  ON c_mv_ec.iso_code  = ec_mv.internal_id
		LEFT JOIN currencies c_mv_iso ON c_mv_iso.iso_char_code = norm_ccy.currency
cross apply dbo.fnFxRateCross(ISNULL(COALESCE(c_mv_ec.iso_char_code,c_mv_iso.iso_char_code), ''),@p2,@p1) cr
		where c.settle_code=c.settle_max`


