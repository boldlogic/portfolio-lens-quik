-- +goose Up
-- +goose StatementBegin
CREATE TABLE moex.markets (
    market_id INT IDENTITY(1,1) PRIMARY KEY,
    market NVARCHAR(50) NOT NULL,
    engine_id INT NOT NULL,
    CONSTRAINT FK_markets_engines FOREIGN KEY (engine_id) REFERENCES moex.engines(engine_id)
);
CREATE UNIQUE INDEX UX_markets_engine_market ON moex.markets(engine_id, market);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS moex.markets;
