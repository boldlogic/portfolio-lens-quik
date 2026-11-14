-- +goose Up

ALTER TABLE quik.trades ALTER COLUMN trade_time datetime NULL;
ALTER TABLE quik.trades ALTER COLUMN trade_time_mcs float NULL;

EXEC sys.sp_updateextendedproperty
    @name = N'MS_Description',
    @value = N'Время регистрации сделки с точностью до секунды (QUIK: «Время»). ODBC QUIK сопоставляет это поле с datetime. Дата внутри значения может быть служебной, дата сделки лежит в trade_date.',
    @level0type = N'SCHEMA', @level0name = N'quik',
    @level1type = N'TABLE', @level1name = N'trades',
    @level2type = N'COLUMN', @level2name = N'trade_time';

EXEC sys.sp_updateextendedproperty
    @name = N'MS_Description',
    @value = N'Микросекунды времени регистрации сделки (QUIK: «Время (мкс)»). ODBC QUIK отдает это число как float.',
    @level0type = N'SCHEMA', @level0name = N'quik',
    @level1type = N'TABLE', @level1name = N'trades',
    @level2type = N'COLUMN', @level2name = N'trade_time_mcs';

-- +goose Down

ALTER TABLE quik.trades ALTER COLUMN trade_time time(0) NULL;
ALTER TABLE quik.trades ALTER COLUMN trade_time_mcs int NULL;

EXEC sys.sp_updateextendedproperty
    @name = N'MS_Description',
    @value = N'Время регистрации сделки с точностью до секунды (QUIK: «Время»). Хранится как его отдал терминал.',
    @level0type = N'SCHEMA', @level0name = N'quik',
    @level1type = N'TABLE', @level1name = N'trades',
    @level2type = N'COLUMN', @level2name = N'trade_time';

EXEC sys.sp_updateextendedproperty
    @name = N'MS_Description',
    @value = N'Микросекунды времени регистрации сделки (QUIK: «Время (мкс)»).',
    @level0type = N'SCHEMA', @level0name = N'quik',
    @level1type = N'TABLE', @level1name = N'trades',
    @level2type = N'COLUMN', @level2name = N'trade_time_mcs';
