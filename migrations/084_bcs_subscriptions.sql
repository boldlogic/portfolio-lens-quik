-- +goose Up

CREATE TABLE bcs.subscriptions (
    subscription_id BIGINT       NOT NULL IDENTITY(1, 1),
    ticker          NVARCHAR(64) NOT NULL,
    class_code      NVARCHAR(32) NOT NULL CONSTRAINT DF_bcs_subscriptions_class DEFAULT N'',
    CONSTRAINT PK_bcs_subscriptions PRIMARY KEY CLUSTERED (subscription_id),
    CONSTRAINT UQ_bcs_subscriptions_ticker_class UNIQUE (ticker, class_code)
);

INSERT INTO bcs.subscriptions (ticker, class_code)
SELECT DISTINCT t.ticker, ISNULL(t.class_code, N'')
FROM bcs.own_trades AS t
WHERE t.ticker <> N''
  AND NOT EXISTS (
      SELECT 1 FROM bcs.subscriptions AS s
      WHERE s.ticker = t.ticker AND s.class_code = ISNULL(t.class_code, N'')
  );

INSERT INTO bcs.subscriptions (ticker, class_code)
SELECT DISTINCT o.ticker, ISNULL(o.class_code, N'')
FROM bcs.cash_operations AS o
WHERE o.ticker IS NOT NULL AND o.ticker <> N''
  AND NOT EXISTS (
      SELECT 1 FROM bcs.subscriptions AS s
      WHERE s.ticker = o.ticker AND s.class_code = ISNULL(o.class_code, N'')
  );

INSERT INTO bcs.subscriptions (ticker, class_code)
SELECT DISTINCT o.isin, N''
FROM bcs.cash_operations AS o
WHERE (o.ticker IS NULL OR o.ticker = N'')
  AND o.isin IS NOT NULL AND o.isin <> N''
  AND NOT EXISTS (
      SELECT 1 FROM bcs.subscriptions AS s
      WHERE s.ticker = o.isin AND s.class_code = N''
  );

INSERT INTO bcs.subscriptions (ticker, class_code)
SELECT DISTINCT sec.ticker, p.board_code
FROM bcs.positions AS p
JOIN bcs.securities AS sec ON sec.security_id = p.security_id
WHERE p.board_code IS NOT NULL AND p.board_code <> N'' AND p.board_code <> N'BQMEA'
  AND sec.ticker <> N''
  AND NOT EXISTS (
      SELECT 1 FROM bcs.subscriptions AS s
      WHERE s.ticker = sec.ticker AND s.class_code = p.board_code
  );

INSERT INTO bcs.subscriptions (ticker, class_code)
SELECT DISTINCT sec.display_name, N''
FROM bcs.positions AS p
JOIN bcs.securities AS sec ON sec.security_id = p.security_id
WHERE p.board_code = N'BQMEA'
  AND sec.display_name IS NOT NULL AND sec.display_name <> N'' AND sec.display_name <> sec.ticker
  AND LEN(sec.display_name) <= 64
  AND NOT EXISTS (
      SELECT 1 FROM bcs.subscriptions AS s
      WHERE s.ticker = sec.display_name AND s.class_code = N''
  );

-- +goose Down

DROP TABLE bcs.subscriptions;
