-- +goose Up
-- +goose StatementBegin
CREATE TABLE moex.historical_rates (
    trade_date DATE NOT NULL,
    board_id NVARCHAR(10) NOT NULL,
    sec_id NVARCHAR(20) NOT NULL,
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
    currency_id NVARCHAR(10) NULL,
    trend_cls_pr FLOAT NULL,
    CONSTRAINT PK_historical_rates PRIMARY KEY (trade_date, board_id, sec_id)
);
CREATE INDEX IX_historical_rates_sec_id ON moex.historical_rates(sec_id);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS moex.historical_rates;
