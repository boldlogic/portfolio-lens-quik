-- +goose Up
-- +goose StatementBegin
CREATE TABLE moex.sync_board_profiles (
    board_id    NVARCHAR(10) NOT NULL PRIMARY KEY,
    engine      NVARCHAR(50) NOT NULL,
    market      NVARCHAR(50) NOT NULL,
    quote_style NVARCHAR(20) NULL,
    enabled     BIT NOT NULL CONSTRAINT DF_sync_board_profiles_enabled DEFAULT 1,
    CONSTRAINT FK_sync_board_profiles_boards FOREIGN KEY (board_id) REFERENCES moex.boards(board_id)
);

CREATE TABLE moex.sync_watchlist (
    watchlist_id    BIGINT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    sec_id          NVARCHAR(20) NOT NULL,
    board_id        NVARCHAR(10) NOT NULL,
    sync_securities BIT NOT NULL CONSTRAINT DF_sync_watchlist_securities DEFAULT 1,
    sync_quotes     BIT NOT NULL CONSTRAINT DF_sync_watchlist_quotes DEFAULT 1,
    sync_history    BIT NOT NULL CONSTRAINT DF_sync_watchlist_history DEFAULT 0,
    enabled         BIT NOT NULL CONSTRAINT DF_sync_watchlist_enabled DEFAULT 1,
    created_at      DATETIME2 NOT NULL CONSTRAINT DF_sync_watchlist_created DEFAULT SYSUTCDATETIME(),
    updated_at      DATETIME2 NULL,
    CONSTRAINT FK_sync_watchlist_boards FOREIGN KEY (board_id) REFERENCES moex.boards(board_id),
    CONSTRAINT UQ_sync_watchlist_sec_board UNIQUE (sec_id, board_id)
);
CREATE INDEX IX_sync_watchlist_board_id ON moex.sync_watchlist(board_id);

CREATE TABLE moex.sync_index_subscriptions (
    index_id          NVARCHAR(20) NOT NULL PRIMARY KEY,
    trade_date        DATE NULL,
    enabled           BIT NOT NULL CONSTRAINT DF_sync_index_subscriptions_enabled DEFAULT 1,
    expand_watchlist  BIT NOT NULL CONSTRAINT DF_sync_index_subscriptions_expand DEFAULT 0
);

-- Seed engines/markets/boards for sync scope (idempotent)
IF NOT EXISTS (SELECT 1 FROM moex.engines WHERE engine = N'stock')
    INSERT INTO moex.engines (engine) VALUES (N'stock');
IF NOT EXISTS (SELECT 1 FROM moex.engines WHERE engine = N'currency')
    INSERT INTO moex.engines (engine) VALUES (N'currency');
IF NOT EXISTS (SELECT 1 FROM moex.engines WHERE engine = N'otc')
    INSERT INTO moex.engines (engine) VALUES (N'otc');

DECLARE @stock_id INT = (SELECT engine_id FROM moex.engines WHERE engine = N'stock');
DECLARE @currency_id INT = (SELECT engine_id FROM moex.engines WHERE engine = N'currency');
DECLARE @otc_id INT = (SELECT engine_id FROM moex.engines WHERE engine = N'otc');

IF NOT EXISTS (SELECT 1 FROM moex.markets WHERE engine_id = @stock_id AND market = N'shares')
    INSERT INTO moex.markets (market, engine_id) VALUES (N'shares', @stock_id);
IF NOT EXISTS (SELECT 1 FROM moex.markets WHERE engine_id = @stock_id AND market = N'bonds')
    INSERT INTO moex.markets (market, engine_id) VALUES (N'bonds', @stock_id);
IF NOT EXISTS (SELECT 1 FROM moex.markets WHERE engine_id = @currency_id AND market = N'selt')
    INSERT INTO moex.markets (market, engine_id) VALUES (N'selt', @currency_id);
IF NOT EXISTS (SELECT 1 FROM moex.markets WHERE engine_id = @currency_id AND market = N'otcindices')
    INSERT INTO moex.markets (market, engine_id) VALUES (N'otcindices', @currency_id);
IF NOT EXISTS (SELECT 1 FROM moex.markets WHERE engine_id = @otc_id AND market = N'shares')
    INSERT INTO moex.markets (market, engine_id) VALUES (N'shares', @otc_id);

DECLARE @m_shares INT = (SELECT market_id FROM moex.markets WHERE engine_id = @stock_id AND market = N'shares');
DECLARE @m_bonds INT = (SELECT market_id FROM moex.markets WHERE engine_id = @stock_id AND market = N'bonds');
DECLARE @m_selt INT = (SELECT market_id FROM moex.markets WHERE engine_id = @currency_id AND market = N'selt');
DECLARE @m_otcidx INT = (SELECT market_id FROM moex.markets WHERE engine_id = @currency_id AND market = N'otcindices');
DECLARE @m_otc_sh INT = (SELECT market_id FROM moex.markets WHERE engine_id = @otc_id AND market = N'shares');

MERGE moex.boards AS t
USING (VALUES
    (N'TQBR', N'TQBR', @m_shares),
    (N'CETS', N'CETS', @m_selt),
    (N'TQCB', N'TQCB', @m_bonds),
    (N'TQTF', N'TQTF', @m_shares),
    (N'COTR', N'COTR', @m_otcidx),
    (N'MTQR', N'MTQR', @m_otc_sh)
) AS s(board_id, title, market_id)
ON t.board_id = s.board_id
WHEN NOT MATCHED THEN
    INSERT (board_id, title, market_id) VALUES (s.board_id, s.title, s.market_id);

MERGE moex.sync_board_profiles AS t
USING (VALUES
    (N'TQBR', N'stock', N'shares', NULL),
    (N'CETS', N'currency', N'selt', NULL),
    (N'TQCB', N'stock', N'bonds', NULL),
    (N'TQTF', N'stock', N'shares', NULL),
    (N'COTR', N'currency', N'otcindices', N'otc_index'),
    (N'MTQR', N'otc', N'shares', NULL)
) AS s(board_id, engine, market, quote_style)
ON t.board_id = s.board_id
WHEN NOT MATCHED THEN
    INSERT (board_id, engine, market, quote_style) VALUES (s.board_id, s.engine, s.market, s.quote_style);

MERGE moex.sync_watchlist AS t
USING (VALUES
    (N'SBER', N'TQBR', 1, 1, 0),
    (N'GAZP', N'TQBR', 1, 1, 0),
    (N'LKOH', N'TQBR', 1, 1, 0),
    (N'VTBR', N'TQBR', 1, 1, 0),
    (N'USD000UTSTOM', N'CETS', 1, 1, 0),
    (N'EUR_RUB__TOM', N'CETS', 1, 1, 0),
    (N'CNYRUB_TOM', N'CETS', 1, 1, 0),
    (N'RU000A105JT4', N'TQCB', 1, 1, 0),
    (N'LQDT', N'TQTF', 1, 1, 0),
    (N'USDRUBTOMOTC', N'COTR', 1, 1, 0),
    (N'EURRUBTOMOTC', N'COTR', 1, 1, 0),
    (N'CNYRUBTOMOTC', N'COTR', 1, 1, 0)
) AS s(sec_id, board_id, sync_securities, sync_quotes, sync_history)
ON t.sec_id = s.sec_id AND t.board_id = s.board_id
WHEN NOT MATCHED THEN
    INSERT (sec_id, board_id, sync_securities, sync_quotes, sync_history)
    VALUES (s.sec_id, s.board_id, s.sync_securities, s.sync_quotes, s.sync_history);

IF NOT EXISTS (SELECT 1 FROM moex.sync_index_subscriptions WHERE index_id = N'IMOEX')
    INSERT INTO moex.sync_index_subscriptions (index_id, trade_date, enabled, expand_watchlist)
    VALUES (N'IMOEX', CAST('2024-12-19' AS DATE), 1, 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS moex.sync_index_subscriptions;
DROP TABLE IF EXISTS moex.sync_watchlist;
DROP TABLE IF EXISTS moex.sync_board_profiles;
-- +goose StatementEnd
