-- +goose Up
-- +goose StatementBegin
CREATE TABLE moex.quotes (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    sec_id NVARCHAR(20) NOT NULL,
    board_id NVARCHAR(10) NOT NULL,
    snapshot_at DATETIME2 NOT NULL DEFAULT GETUTCDATE(),
    last_price DECIMAL(18,6),
    open_price DECIMAL(18,6),
    high_price DECIMAL(18,6),
    low_price DECIMAL(18,6),
    close_price DECIMAL(18,6),
    bid_price DECIMAL(18,6),
    offer_price DECIMAL(18,6),
    bid_depth DECIMAL(18,2),
    offer_depth DECIMAL(18,2),
    volume BIGINT,
    value DECIMAL(18,2),
    num_trades INT,
    waprice DECIMAL(18,6),
    market_price DECIMAL(18,6),
    last_change DECIMAL(18,6),
    last_change_prc DECIMAL(10,4),
    accrued_int DECIMAL(18,6),
    yield_to_mat DECIMAL(10,4),
    duration DECIMAL(10,4),
    CONSTRAINT FK_quotes_securities FOREIGN KEY (sec_id) REFERENCES moex.securities(sec_id),
    CONSTRAINT FK_quotes_boards FOREIGN KEY (board_id) REFERENCES moex.boards(board_id)
);
CREATE INDEX IX_quotes_sec_board ON moex.quotes(sec_id, board_id);
CREATE INDEX IX_quotes_snapshot_at ON moex.quotes(snapshot_at);
CREATE INDEX IX_quotes_sec_snapshot ON moex.quotes(sec_id, snapshot_at DESC);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS moex.quotes;
