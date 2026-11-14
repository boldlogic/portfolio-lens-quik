-- +goose Up

ALTER TABLE bcs.own_trades ADD isin NVARCHAR(16) NULL;
ALTER TABLE bcs.own_trades ALTER COLUMN ticker NVARCHAR(64) NULL;

UPDATE t
SET t.security_id = COALESCE(t.security_id, s.security_id),
    t.ticker = s.ticker,
    t.isin = s.isin
FROM bcs.own_trades AS t
JOIN bcs.securities AS s ON s.isin = t.ticker
WHERE t.isin IS NULL
  AND UPPER(t.ticker) COLLATE Latin1_General_100_BIN2 LIKE
      '[A-Z][A-Z][A-Z0-9][A-Z0-9][A-Z0-9][A-Z0-9][A-Z0-9][A-Z0-9][A-Z0-9][A-Z0-9][A-Z0-9][0-9]'
  AND (SELECT COUNT(*) FROM bcs.securities WHERE isin = t.ticker) = 1;

UPDATE t
SET t.isin = UPPER(t.ticker),
    t.ticker = NULL
FROM bcs.own_trades AS t
WHERE t.isin IS NULL
  AND UPPER(t.ticker) COLLATE Latin1_General_100_BIN2 LIKE
      '[A-Z][A-Z][A-Z0-9][A-Z0-9][A-Z0-9][A-Z0-9][A-Z0-9][A-Z0-9][A-Z0-9][A-Z0-9][A-Z0-9][0-9]'
  AND NOT EXISTS (SELECT 1 FROM bcs.securities WHERE ticker = t.ticker);

-- +goose Down

IF EXISTS (SELECT 1 FROM bcs.own_trades WHERE ticker IS NULL)
    THROW 51000, 'Cannot remove ISIN column while trades rely on it without a ticker.', 1;

ALTER TABLE bcs.own_trades ALTER COLUMN ticker NVARCHAR(64) NOT NULL;
ALTER TABLE bcs.own_trades DROP COLUMN isin;
