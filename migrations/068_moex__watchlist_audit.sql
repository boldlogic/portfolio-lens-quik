-- +goose Up
-- +goose StatementBegin
ALTER TABLE moex.sync_watchlist ADD
    source       NVARCHAR(50) NULL,
    requested_by NVARCHAR(100) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE moex.sync_watchlist DROP COLUMN source, requested_by;
-- +goose StatementEnd
