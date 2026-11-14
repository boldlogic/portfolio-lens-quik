-- +goose Up
-- +goose StatementBegin
SET XACT_ABORT ON;
-- +goose StatementEnd

-- +goose StatementBegin

IF OBJECT_ID(N'tinvest.history_record', N'U') IS NULL
    THROW 51000, N'Отсутствует tinvest.history_record: перенос истории невозможен.', 1;
IF OBJECT_ID(N'tinvest.operation', N'U') IS NULL OR OBJECT_ID(N'tinvest.insider_deal', N'U') IS NULL
    THROW 51000, N'Не найдены типизированные таблицы T-Invest.', 1;

IF EXISTS (
    SELECT 1
    FROM tinvest.history_record h
    CROSS APPLY OPENJSON(h.payload) WITH (
        executed_at NVARCHAR(80) '$.at',
        payment_amount NVARCHAR(128) '$.payment.amount',
        quantity NVARCHAR(64) '$.quantity'
    ) j
    WHERE h.entity = N'operation'
      AND (TRY_CONVERT(DATETIMEOFFSET(7), j.executed_at) IS NULL
        OR TRY_CONVERT(DECIMAL(28,9), j.payment_amount) IS NULL
        OR TRY_CONVERT(DECIMAL(28,9), j.quantity) IS NULL)
)
    THROW 51000, N'В JSON операций найдены значения, которые нельзя перенести в типизированные столбцы.', 1;

IF EXISTS (
    SELECT 1
    FROM tinvest.history_record h
    CROSS APPLY OPENJSON(h.payload,N'$.fills') f
    WHERE h.entity=N'operation'
      AND (TRY_CONVERT(DECIMAL(28,9),JSON_VALUE(f.value,N'$.price.amount')) IS NULL
        OR TRY_CONVERT(DECIMAL(28,9),JSON_VALUE(f.value,N'$.quantity')) IS NULL
        OR (NULLIF(JSON_VALUE(f.value,N'$.at'),N'') IS NOT NULL
            AND TRY_CONVERT(DATETIMEOFFSET(7),JSON_VALUE(f.value,N'$.at')) IS NULL))
)
    THROW 51000, N'В fills операций найдены значения, которые нельзя перенести в типизированные столбцы.', 1;

IF EXISTS (
    SELECT 1
    FROM tinvest.history_record h
    CROSS APPLY OPENJSON(h.payload) WITH (
        executed_at NVARCHAR(80) '$.at',
        quantity NVARCHAR(64) '$.quantity',
        price_amount NVARCHAR(128) '$.price.amount',
        percentage NVARCHAR(128) '$.percentage'
    ) j
    WHERE h.entity = N'insider_deal'
      AND (TRY_CONVERT(DECIMAL(28,9), j.quantity) IS NULL
        OR TRY_CONVERT(DECIMAL(28,9), j.price_amount) IS NULL
        OR (NULLIF(j.executed_at, N'') IS NOT NULL AND TRY_CONVERT(DATETIMEOFFSET(7), j.executed_at) IS NULL)
        OR (NULLIF(j.percentage, N'') IS NOT NULL AND TRY_CONVERT(DECIMAL(28,9), j.percentage) IS NULL))
)
    THROW 51000, N'В JSON инсайдерских сделок найдены значения, которые нельзя перенести в типизированные столбцы.', 1;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE tinvest.operation ADD operation_id BIGINT IDENTITY(1,1) NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE tinvest.operation ADD CONSTRAINT PK_tinvest_operation PRIMARY KEY CLUSTERED (operation_id);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE tinvest.operation ADD operation_key CHAR(64) NULL;
-- +goose StatementEnd

-- +goose StatementBegin
DECLARE @constraint_name SYSNAME;
DECLARE @drop_constraint_sql NVARCHAR(4000);
DECLARE unique_constraints CURSOR LOCAL FAST_FORWARD FOR
    SELECT kc.name
    FROM sys.key_constraints kc
    WHERE kc.parent_object_id = OBJECT_ID(N'tinvest.operation') AND kc.type = N'UQ';
OPEN unique_constraints;
FETCH NEXT FROM unique_constraints INTO @constraint_name;
WHILE @@FETCH_STATUS = 0
BEGIN
    SET @drop_constraint_sql=N'ALTER TABLE tinvest.operation DROP CONSTRAINT ' + QUOTENAME(@constraint_name);
    EXEC(@drop_constraint_sql);
    FETCH NEXT FROM unique_constraints INTO @constraint_name;
END;
CLOSE unique_constraints;
DEALLOCATE unique_constraints;
-- +goose StatementEnd

-- +goose StatementBegin
IF EXISTS (
    SELECT 1 FROM tinvest.operation
    WHERE TRY_CONVERT(DECIMAL(28,9), payment_amount) IS NULL
)
    THROW 51000, N'В tinvest.operation найдены суммы, которые нельзя преобразовать в DECIMAL(28,9).', 1;
ALTER TABLE tinvest.operation ALTER COLUMN payment_amount DECIMAL(28,9) NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE tinvest.operation_fill (
    operation_id BIGINT NOT NULL,
    fill_index INT NOT NULL,
    trade_num NVARCHAR(64) NOT NULL DEFAULT N'',
    executed_at DATETIMEOFFSET(7) NULL,
    instrument_uid NVARCHAR(64) NOT NULL DEFAULT N'',
    figi NVARCHAR(32) NOT NULL DEFAULT N'',
    ticker NVARCHAR(64) NOT NULL DEFAULT N'',
    account_id NVARCHAR(128) NOT NULL DEFAULT N'',
    direction NVARCHAR(128) NOT NULL DEFAULT N'',
    price_amount DECIMAL(28,9) NOT NULL,
    price_currency NVARCHAR(16) NOT NULL DEFAULT N'',
    quantity DECIMAL(28,9) NOT NULL,
    CONSTRAINT PK_tinvest_operation_fill PRIMARY KEY (operation_id, fill_index),
    CONSTRAINT FK_tinvest_operation_fill_operation FOREIGN KEY (operation_id)
        REFERENCES tinvest.operation(operation_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE tinvest.insider_deal ADD record_key CHAR(64) NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX UQ_tinvest_operation_natural
    ON tinvest.operation(account_id, executed_at, type, instrument_uid, payment_amount, payment_currency, quantity);
CREATE UNIQUE INDEX UQ_tinvest_operation_key ON tinvest.operation(operation_key) WHERE operation_key IS NOT NULL;
CREATE UNIQUE INDEX UQ_tinvest_insider_record_key ON tinvest.insider_deal(record_key) WHERE record_key IS NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE #operations (
    scope_key NVARCHAR(128) NOT NULL,
    record_key CHAR(64) NOT NULL,
    executed_at DATETIMEOFFSET(7) NOT NULL,
    type NVARCHAR(64) NOT NULL,
    instrument_uid NVARCHAR(64) NOT NULL,
    payment_amount DECIMAL(28,9) NOT NULL,
    payment_currency NVARCHAR(16) NOT NULL,
    quantity DECIMAL(28,9) NOT NULL,
    name NVARCHAR(200) NOT NULL,
    state NVARCHAR(64) NOT NULL,
    figi NVARCHAR(32) NOT NULL,
    ticker NVARCHAR(64) NOT NULL,
    PRIMARY KEY (scope_key, record_key)
);

INSERT INTO #operations(scope_key,record_key,executed_at,type,instrument_uid,payment_amount,payment_currency,quantity,name,state,figi,ticker)
SELECT h.scope_key,h.record_key,CONVERT(DATETIMEOFFSET(7),j.executed_at),ISNULL(j.type,N''),ISNULL(j.instrument_uid,N''),
       CONVERT(DECIMAL(28,9),j.payment_amount),ISNULL(j.payment_currency,N''),CONVERT(DECIMAL(28,9),j.quantity),
       ISNULL(j.name,N''),ISNULL(j.state,N''),ISNULL(j.figi,N''),ISNULL(j.ticker,N'')
FROM tinvest.history_record h
CROSS APPLY OPENJSON(h.payload) WITH (
    executed_at DATETIMEOFFSET(7) '$.at', type NVARCHAR(64) '$.type', name NVARCHAR(200) '$.name',
    state NVARCHAR(64) '$.state', figi NVARCHAR(32) '$.figi', ticker NVARCHAR(64) '$.ticker',
    instrument_uid NVARCHAR(64) '$.instrument_uid', payment_amount DECIMAL(28,9) '$.payment.amount',
    payment_currency NVARCHAR(16) '$.payment.currency', quantity DECIMAL(28,9) '$.quantity'
) j
WHERE h.entity=N'operation';
-- +goose StatementEnd

-- +goose StatementBegin
DECLARE @op_scope NVARCHAR(128), @op_key CHAR(64), @op_at DATETIMEOFFSET(7), @op_type NVARCHAR(64),
        @op_uid NVARCHAR(64), @op_amount DECIMAL(28,9), @op_currency NVARCHAR(16), @op_quantity DECIMAL(28,9),
        @op_name NVARCHAR(200), @op_state NVARCHAR(64), @op_figi NVARCHAR(32), @op_ticker NVARCHAR(64);
DECLARE operation_rows CURSOR LOCAL FAST_FORWARD FOR
    SELECT scope_key,record_key,executed_at,type,instrument_uid,payment_amount,payment_currency,quantity,name,state,figi,ticker
    FROM #operations ORDER BY scope_key,record_key;
OPEN operation_rows;
FETCH NEXT FROM operation_rows INTO @op_scope,@op_key,@op_at,@op_type,@op_uid,@op_amount,@op_currency,@op_quantity,@op_name,@op_state,@op_figi,@op_ticker;
WHILE @@FETCH_STATUS=0
BEGIN
    UPDATE tinvest.operation
    SET operation_key=@op_key,name=@op_name,state=@op_state,figi=@op_figi,ticker=@op_ticker
    WHERE operation_key=@op_key
       OR (account_id=@op_scope AND executed_at=@op_at AND type=@op_type AND instrument_uid=@op_uid
       AND payment_amount=@op_amount AND payment_currency=@op_currency AND quantity=@op_quantity);
    IF @@ROWCOUNT=0
        INSERT INTO tinvest.operation(account_id,executed_at,type,instrument_uid,payment_amount,payment_currency,quantity,name,state,figi,ticker,operation_key)
        VALUES(@op_scope,@op_at,@op_type,@op_uid,@op_amount,@op_currency,@op_quantity,@op_name,@op_state,@op_figi,@op_ticker,@op_key);
    FETCH NEXT FROM operation_rows INTO @op_scope,@op_key,@op_at,@op_type,@op_uid,@op_amount,@op_currency,@op_quantity,@op_name,@op_state,@op_figi,@op_ticker;
END;
CLOSE operation_rows;
DEALLOCATE operation_rows;

INSERT INTO tinvest.operation_fill(operation_id,fill_index,trade_num,executed_at,instrument_uid,figi,ticker,account_id,direction,price_amount,price_currency,quantity)
SELECT operation.operation_id,CONVERT(INT,fill_json.[key]),ISNULL(fill.trade_num,N''),TRY_CONVERT(DATETIMEOFFSET(7),fill.executed_at),
       ISNULL(fill.instrument_uid,N''),ISNULL(fill.figi,N''),ISNULL(fill.ticker,N''),ISNULL(fill.account_id,N''),ISNULL(fill.direction,N''),
       TRY_CONVERT(DECIMAL(28,9),fill.price_amount),ISNULL(fill.price_currency,N''),TRY_CONVERT(DECIMAL(28,9),fill.quantity)
FROM #operations source
JOIN tinvest.history_record history ON history.entity=N'operation' AND history.scope_key=source.scope_key AND history.record_key=source.record_key
JOIN tinvest.operation operation ON operation.operation_key=source.record_key
CROSS APPLY OPENJSON(history.payload,N'$.fills') fill_json
CROSS APPLY OPENJSON(fill_json.value) WITH (
    trade_num NVARCHAR(64) '$.trade_num', executed_at NVARCHAR(80) '$.at', instrument_uid NVARCHAR(64) '$.instrument_uid',
    figi NVARCHAR(32) '$.figi', ticker NVARCHAR(64) '$.ticker', account_id NVARCHAR(128) '$.account_id',
    direction NVARCHAR(128) '$.direction', price_amount NVARCHAR(128) '$.price.amount',
    price_currency NVARCHAR(16) '$.price.currency', quantity NVARCHAR(64) '$.quantity'
) fill
WHERE TRY_CONVERT(DECIMAL(28,9),fill.price_amount) IS NOT NULL
  AND TRY_CONVERT(DECIMAL(28,9),fill.quantity) IS NOT NULL
OPTION (LOOP JOIN, MAXDOP 1, MAX_GRANT_PERCENT=1);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE #insider_deals (
    record_key CHAR(64) NOT NULL PRIMARY KEY,
    trade_id BIGINT NOT NULL,
    instrument_uid NVARCHAR(64) NOT NULL,
    direction NVARCHAR(128) NOT NULL,
    currency NVARCHAR(16) NOT NULL,
    executed_at DATETIMEOFFSET(7) NULL,
    quantity DECIMAL(28,9) NOT NULL,
    price_amount DECIMAL(28,9) NOT NULL,
    ticker NVARCHAR(64) NOT NULL,
    investor_name NVARCHAR(200) NOT NULL,
    investor_position NVARCHAR(200) NOT NULL,
    percentage DECIMAL(28,9) NULL,
    option_execution BIT NOT NULL,
    disclosure_at DATETIMEOFFSET(7) NULL
);

INSERT INTO #insider_deals(record_key,trade_id,instrument_uid,direction,currency,executed_at,quantity,price_amount,ticker,investor_name,investor_position,percentage,option_execution,disclosure_at)
SELECT h.record_key,ISNULL(j.trade_id,0),ISNULL(j.instrument_uid,N''),ISNULL(j.direction,N''),ISNULL(j.currency,N''),
       TRY_CONVERT(DATETIMEOFFSET(7),j.executed_at),TRY_CONVERT(DECIMAL(28,9),j.quantity),TRY_CONVERT(DECIMAL(28,9),j.price_amount),
       ISNULL(j.ticker,N''),ISNULL(j.investor_name,N''),ISNULL(j.investor_position,N''),TRY_CONVERT(DECIMAL(28,9),NULLIF(j.percentage,N'')),
       CASE WHEN LOWER(j.option_execution) IN (N'true',N'1') THEN 1 ELSE 0 END,TRY_CONVERT(DATETIMEOFFSET(7),NULLIF(j.disclosure_at,N''))
FROM tinvest.history_record h
CROSS APPLY OPENJSON(h.payload) WITH (
    trade_id BIGINT '$.trade_id', instrument_uid NVARCHAR(64) '$.instrument_uid', direction NVARCHAR(128) '$.direction',
    currency NVARCHAR(16) '$.currency', executed_at NVARCHAR(80) '$.at', quantity NVARCHAR(64) '$.quantity',
    price_amount NVARCHAR(128) '$.price.amount', ticker NVARCHAR(64) '$.ticker', investor_name NVARCHAR(200) '$.investor_name',
    investor_position NVARCHAR(200) '$.investor_position', percentage NVARCHAR(128) '$.percentage',
    option_execution NVARCHAR(8) '$.option_execution', disclosure_at NVARCHAR(80) '$.disclosure'
) j
WHERE h.entity=N'insider_deal';
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE target
SET record_key=source.record_key,instrument_uid=source.instrument_uid,direction=source.direction,currency=source.currency,
    executed_at=source.executed_at,quantity=source.quantity,price_amount=source.price_amount,ticker=source.ticker,
    investor_name=source.investor_name,investor_position=source.investor_position,percentage=source.percentage,
    option_execution=source.option_execution,disclosure_at=source.disclosure_at
FROM tinvest.insider_deal target
JOIN #insider_deals source
  ON (source.trade_id<>0 AND target.trade_id=source.trade_id)
  OR (source.trade_id=0 AND target.trade_id=0 AND target.instrument_uid=source.instrument_uid
      AND (target.executed_at=source.executed_at OR (target.executed_at IS NULL AND source.executed_at IS NULL))
      AND target.investor_name=source.investor_name AND target.direction=source.direction
      AND target.quantity=source.quantity AND target.price_amount=source.price_amount);

INSERT INTO tinvest.insider_deal(record_key,trade_id,instrument_uid,direction,currency,executed_at,quantity,price_amount,ticker,investor_name,investor_position,percentage,option_execution,disclosure_at)
SELECT source.record_key,source.trade_id,source.instrument_uid,source.direction,source.currency,source.executed_at,source.quantity,
       source.price_amount,source.ticker,source.investor_name,source.investor_position,source.percentage,source.option_execution,source.disclosure_at
FROM #insider_deals source
WHERE NOT EXISTS (
    SELECT 1 FROM tinvest.insider_deal target
    WHERE target.record_key=source.record_key
       OR (source.trade_id<>0 AND target.trade_id=source.trade_id)
       OR (source.trade_id=0 AND target.trade_id=0 AND target.instrument_uid=source.instrument_uid
           AND (target.executed_at=source.executed_at OR (target.executed_at IS NULL AND source.executed_at IS NULL))
           AND target.investor_name=source.investor_name AND target.direction=source.direction
           AND target.quantity=source.quantity AND target.price_amount=source.price_amount)
);

DELETE FROM tinvest.history_record WHERE entity IN (N'operation',N'insider_deal');

GRANT SELECT, INSERT, UPDATE ON OBJECT::tinvest.operation TO t_invest_portfolio_writer;
GRANT REFERENCES ON OBJECT::tinvest.operation TO t_invest_portfolio_writer;
GRANT SELECT, INSERT, DELETE ON OBJECT::tinvest.operation_fill TO t_invest_portfolio_writer;
GRANT SELECT, INSERT, UPDATE ON OBJECT::tinvest.insider_deal TO t_invest_portfolio_writer;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
THROW 51000, N'Миграция 076 переносит историю в нормализованные таблицы и не откатывается автоматически.', 1;
-- +goose StatementEnd
