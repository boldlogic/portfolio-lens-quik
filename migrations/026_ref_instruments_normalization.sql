-- +goose Up

IF EXISTS (
    SELECT 1
    FROM sys.indexes
    WHERE name = N'UQ_ref_instruments_ticker_trade_point_id'
        AND object_id = OBJECT_ID(N'ref.instruments')
)
DROP INDEX UQ_ref_instruments_ticker_trade_point_id ON ref.instruments;

IF EXISTS (
    SELECT 1
    FROM sys.foreign_keys
    WHERE name = N'FK_ref_instruments_trade_point_id'
        AND parent_object_id = OBJECT_ID(N'ref.instruments')
)
ALTER TABLE ref.instruments DROP CONSTRAINT FK_ref_instruments_trade_point_id;

IF COL_LENGTH('ref.instruments', 'type_id') IS NULL
ALTER TABLE ref.instruments
    ADD type_id TINYINT NULL;

UPDATE i
SET type_id = src.type_id
FROM ref.instruments i
JOIN (
    SELECT
        ib.instrument_id,
        type_id = MIN(itb.type_id)
    FROM ref.instrument_boards ib
    JOIN ref.instrument_type_boards itb
        ON itb.board_id = ib.board_id
    GROUP BY ib.instrument_id
) src
    ON src.instrument_id = i.instrument_id
WHERE i.type_id IS NULL;

-- +goose StatementBegin
IF EXISTS (
    SELECT 1
    FROM ref.instruments
    WHERE type_id IS NULL
)
BEGIN
    THROW 51000, N'Для существующих инструментов не удалось определить тип', 1;
END
-- +goose StatementEnd

IF EXISTS (
    SELECT 1
    FROM sys.columns
    WHERE object_id = OBJECT_ID(N'ref.instruments')
        AND name = N'type_id'
        AND is_nullable = 1
)
ALTER TABLE ref.instruments
    ALTER COLUMN type_id TINYINT NOT NULL;

IF EXISTS (
    SELECT 1
    FROM sys.columns
    WHERE object_id = OBJECT_ID(N'ref.instruments')
        AND name = N'trade_point_id'
        AND is_nullable = 0
)
ALTER TABLE ref.instruments
    ALTER COLUMN trade_point_id TINYINT NULL;

IF NOT EXISTS (
    SELECT 1
    FROM sys.foreign_keys
    WHERE name = N'FK_ref_instruments_trade_point_id'
        AND parent_object_id = OBJECT_ID(N'ref.instruments')
)
ALTER TABLE ref.instruments
    ADD CONSTRAINT FK_ref_instruments_trade_point_id FOREIGN KEY (trade_point_id)
        REFERENCES ref.trade_points (point_id);

IF NOT EXISTS (
    SELECT 1
    FROM sys.foreign_keys
    WHERE name = N'FK_ref_instruments_type_id'
        AND parent_object_id = OBJECT_ID(N'ref.instruments')
)
ALTER TABLE ref.instruments
    ADD CONSTRAINT FK_ref_instruments_type_id FOREIGN KEY (type_id)
        REFERENCES ref.instrument_types (type_id);

IF NOT EXISTS (
    SELECT 1
    FROM sys.indexes
    WHERE name = N'UQ_ref_instruments_ticker_trade_point_type'
        AND object_id = OBJECT_ID(N'ref.instruments')
)
CREATE UNIQUE NONCLUSTERED INDEX UQ_ref_instruments_ticker_trade_point_type
    ON ref.instruments (ticker, trade_point_id, type_id)
    WHERE trade_point_id IS NOT NULL;

-- +goose StatementBegin
DECLARE @currency_type_id TINYINT = (
    SELECT type_id
    FROM ref.instrument_types
    WHERE title = N'Валюта'
);

IF @currency_type_id IS NULL
BEGIN
    THROW 51001, N'Тип инструмента Валюта не найден', 1;
END;

IF NOT EXISTS (
    SELECT 1
    FROM sys.indexes
    WHERE name = N'UQ_ref_instruments_currency_ticker_type'
        AND object_id = OBJECT_ID(N'ref.instruments')
)
BEGIN
    DECLARE @create_currency_index_sql NVARCHAR(MAX) = N'
CREATE UNIQUE NONCLUSTERED INDEX UQ_ref_instruments_currency_ticker_type
    ON ref.instruments (ticker, type_id)
    WHERE trade_point_id IS NULL
        AND type_id = ' + CAST(@currency_type_id AS NVARCHAR(10)) + N';';

    EXEC (@create_currency_index_sql);
END
-- +goose StatementEnd

IF COL_LENGTH('ref.instrument_boards', 'rw') IS NULL
ALTER TABLE ref.instrument_boards
    ADD rw ROWVERSION NOT NULL;
