-- +goose Up
-- +goose StatementBegin
CREATE TABLE moex.boards (
    board_id NVARCHAR(10) PRIMARY KEY,
    title NVARCHAR(100),
    market_id INT NOT NULL,
    currency_id NVARCHAR(10),
    unit NVARCHAR(10),
    CONSTRAINT FK_boards_markets FOREIGN KEY (market_id) REFERENCES moex.markets(market_id)
);
CREATE INDEX IX_boards_market_id ON moex.boards(market_id);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS moex.boards;
