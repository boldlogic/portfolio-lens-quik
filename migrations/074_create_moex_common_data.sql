-- +goose Up

IF NOT EXISTS (SELECT 1 FROM ref.external_systems WHERE ext_system = 'MOEX')
    INSERT INTO ref.external_systems (ext_system) VALUES ('MOEX');

IF NOT EXISTS (SELECT 1 FROM ref.external_code_types WHERE ext_code_type_id = 3)
    INSERT INTO ref.external_code_types (ext_code_type_id, code, name)
    VALUES (3, 'instrument_code', N'Код инструмента источника');

CREATE TABLE moex.common_board_map (
    moex_board_id NVARCHAR(10) NOT NULL,
    board_id TINYINT NOT NULL,
    CONSTRAINT PK_moex_common_board_map PRIMARY KEY CLUSTERED (moex_board_id),
    CONSTRAINT FK_moex_common_board_map_board FOREIGN KEY (board_id) REFERENCES ref.boards (board_id)
);

INSERT INTO moex.common_board_map (moex_board_id, board_id)
SELECT mb.board_id, rb.board_id
FROM moex.boards mb
JOIN ref.boards rb ON rb.code = mb.board_id
JOIN moex.markets mm ON mm.market_id = mb.market_id
JOIN moex.engines me ON me.engine_id = mm.engine_id
WHERE me.engine IN (N'stock', N'currency', N'futures', N'commodity')
  AND NOT EXISTS (SELECT 1 FROM moex.common_board_map m WHERE m.moex_board_id = mb.board_id);

CREATE TABLE market.moex_quotes (
    instrument_id BIGINT NOT NULL,
    board_id TINYINT NOT NULL,
    quote_date DATE NOT NULL,
    snapshot_at DATETIME2 NOT NULL,
    last_price DECIMAL(18,6) NULL,
    open_price DECIMAL(18,6) NULL,
    high_price DECIMAL(18,6) NULL,
    low_price DECIMAL(18,6) NULL,
    close_price DECIMAL(18,6) NULL,
    bid_price DECIMAL(18,6) NULL,
    offer_price DECIMAL(18,6) NULL,
    bid_depth DECIMAL(18,2) NULL,
    offer_depth DECIMAL(18,2) NULL,
    volume BIGINT NULL,
    value DECIMAL(18,2) NULL,
    num_trades INT NULL,
    waprice DECIMAL(18,6) NULL,
    market_price DECIMAL(18,6) NULL,
    last_change DECIMAL(18,6) NULL,
    last_change_prc DECIMAL(10,4) NULL,
    accrued_int DECIMAL(18,6) NULL,
    yield_to_mat DECIMAL(10,4) NULL,
    duration DECIMAL(10,4) NULL,
    CONSTRAINT PK_market_moex_quotes PRIMARY KEY CLUSTERED (instrument_id, board_id, quote_date),
    CONSTRAINT FK_market_moex_quotes_instrument FOREIGN KEY (instrument_id) REFERENCES ref.instruments (instrument_id),
    CONSTRAINT FK_market_moex_quotes_board FOREIGN KEY (board_id) REFERENCES ref.boards (board_id)
);

CREATE TABLE market.moex_historical_rates (
    instrument_id BIGINT NOT NULL,
    board_id TINYINT NOT NULL,
    trade_date DATE NOT NULL,
    short_name NVARCHAR(100) NULL,
    num_trades FLOAT NULL,
    value FLOAT NULL,
    open_price FLOAT NULL,
    low_price FLOAT NULL,
    high_price FLOAT NULL,
    legal_close_price FLOAT NULL,
    waprice FLOAT NULL,
    close_price FLOAT NULL,
    volume FLOAT NULL,
    market_price2 FLOAT NULL,
    market_price3 FLOAT NULL,
    admitted_quote FLOAT NULL,
    mp2_val_trd FLOAT NULL,
    market_price3_trades_value FLOAT NULL,
    admitted_value FLOAT NULL,
    waval FLOAT NULL,
    trading_session FLOAT NULL,
    currency_id SMALLINT NULL,
    trend_cls_pr FLOAT NULL,
    CONSTRAINT PK_market_moex_historical_rates PRIMARY KEY CLUSTERED (instrument_id, board_id, trade_date),
    CONSTRAINT FK_market_moex_history_instrument FOREIGN KEY (instrument_id) REFERENCES ref.instruments (instrument_id),
    CONSTRAINT FK_market_moex_history_board FOREIGN KEY (board_id) REFERENCES ref.boards (board_id),
    CONSTRAINT FK_market_moex_history_currency FOREIGN KEY (currency_id) REFERENCES ref.currencies (iso_code)
);

CREATE TABLE market.moex_indices (
    index_code NVARCHAR(20) NOT NULL,
    CONSTRAINT PK_market_moex_indices PRIMARY KEY CLUSTERED (index_code)
);

CREATE TABLE market.moex_index_components (
    index_code NVARCHAR(20) NOT NULL,
    trade_date DATE NOT NULL,
    ticker NVARCHAR(20) NOT NULL,
    short_name NVARCHAR(100) NULL,
    sec_id NVARCHAR(20) NULL,
    instrument_id BIGINT NULL,
    weight DECIMAL(10,4) NULL,
    trading_session INT NULL,
    CONSTRAINT PK_market_moex_index_components PRIMARY KEY CLUSTERED (index_code, trade_date, ticker),
    CONSTRAINT FK_market_moex_components_index FOREIGN KEY (index_code) REFERENCES market.moex_indices (index_code),
    CONSTRAINT FK_market_moex_components_instrument FOREIGN KEY (instrument_id) REFERENCES ref.instruments (instrument_id)
);

-- +goose Down
THROW 51000, N'Миграция 074 содержит импортированные данные MOEX и намеренно не откатывается автоматически.', 1;
