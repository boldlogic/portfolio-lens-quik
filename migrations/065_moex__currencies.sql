-- +goose Up
-- +goose StatementBegin
CREATE TABLE moex.currencies (
    currency_id NVARCHAR(10) NOT NULL PRIMARY KEY,
    title NVARCHAR(100) NULL
);

UPDATE moex.securities SET currency_id = NULL WHERE currency_id IS NOT NULL AND LTRIM(RTRIM(currency_id)) = '';
UPDATE moex.boards SET currency_id = NULL WHERE currency_id IS NOT NULL AND LTRIM(RTRIM(currency_id)) = '';
UPDATE moex.historical_rates SET currency_id = NULL WHERE currency_id IS NOT NULL AND LTRIM(RTRIM(currency_id)) = '';

INSERT INTO moex.currencies (currency_id, title)
SELECT DISTINCT v.currency_id, NULL
FROM (
    SELECT NULLIF(LTRIM(RTRIM(currency_id)), '') AS currency_id FROM moex.securities
    UNION
    SELECT NULLIF(LTRIM(RTRIM(currency_id)), '') FROM moex.boards
    UNION
    SELECT NULLIF(LTRIM(RTRIM(currency_id)), '') FROM moex.historical_rates
) AS v
WHERE v.currency_id IS NOT NULL;

INSERT INTO moex.currencies (currency_id, title)
SELECT c.currency_id, c.title
FROM (VALUES
    ('RUB', N'Российский рубль'),
    ('USD', N'Доллар США'),
    ('EUR', N'Евро'),
    ('CNY', N'Китайский юань'),
    ('SUR', N'Рубль (legacy MOEX)'),
    ('GBP', N'Фунт стерлингов'),
    ('HKD', N'Гонконгский доллар'),
    ('TRY', N'Турецкая лира'),
    ('JPY', N'Японская иена')
) AS c(currency_id, title)
WHERE NOT EXISTS (SELECT 1 FROM moex.currencies x WHERE x.currency_id = c.currency_id);

UPDATE moex.currencies SET title = c.title
FROM moex.currencies cur
INNER JOIN (VALUES
    ('RUB', N'Российский рубль'),
    ('USD', N'Доллар США'),
    ('EUR', N'Евро'),
    ('CNY', N'Китайский юань'),
    ('SUR', N'Рубль (legacy MOEX)'),
    ('GBP', N'Фунт стерлингов'),
    ('HKD', N'Гонконгский доллар'),
    ('TRY', N'Турецкая лира'),
    ('JPY', N'Японская иена')
) AS c(currency_id, title) ON cur.currency_id = c.currency_id
WHERE cur.title IS NULL;

ALTER TABLE moex.securities
    ADD CONSTRAINT FK_securities_currencies
    FOREIGN KEY (currency_id) REFERENCES moex.currencies(currency_id);

ALTER TABLE moex.boards
    ADD CONSTRAINT FK_boards_currencies
    FOREIGN KEY (currency_id) REFERENCES moex.currencies(currency_id);

ALTER TABLE moex.historical_rates
    ADD CONSTRAINT FK_historical_rates_currencies
    FOREIGN KEY (currency_id) REFERENCES moex.currencies(currency_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE moex.historical_rates DROP CONSTRAINT FK_historical_rates_currencies;
ALTER TABLE moex.boards DROP CONSTRAINT FK_boards_currencies;
ALTER TABLE moex.securities DROP CONSTRAINT FK_securities_currencies;
DROP TABLE IF EXISTS moex.currencies;
-- +goose StatementEnd
