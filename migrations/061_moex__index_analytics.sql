-- +goose Up
-- +goose StatementBegin
CREATE TABLE moex.index_analytics (
    index_id NVARCHAR(20) NOT NULL,
    trade_date DATE NOT NULL,
    ticker NVARCHAR(20) NOT NULL,
    short_name NVARCHAR(100) NULL,
    sec_id NVARCHAR(20) NULL,
    weight DECIMAL(10,4) NULL,
    trading_session INT NULL,
    CONSTRAINT PK_index_analytics PRIMARY KEY (index_id, trade_date, ticker)
);
CREATE INDEX IX_index_analytics_sec_id ON moex.index_analytics(sec_id);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS moex.index_analytics;
