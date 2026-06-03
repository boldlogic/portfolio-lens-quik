declare @ext_system_id tinyint = (
    select
        ext_system_id
    from
        dbo.external_systems
    where
        ext_system = 'QUIK'
);

declare @load_date date = '2026-06-03';

-- WITH cte AS (
    SELECT
        li.load_date,
        li.source_date,
        li.client_code,
        currency_code = case when UPPER(TRIM(li.ccy)) in ('SUR', 'RUR') THEN 'RUB' ELSE UPPER(TRIM(li.ccy)) END,
        li.settle_code,
        li.firm_code,
        li.firm_name,
        li.balance,
        settle_max = MAX(li.settle_code) OVER (
            PARTITION BY li.load_date,
            li.client_code,
            li.ccy,
            li.position_code,
            li.firm_code
        )
    FROM
        quik.money_limits li
    WHERE
        li.load_date = @load_date