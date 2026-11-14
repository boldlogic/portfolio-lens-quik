-- +goose Up
-- +goose StatementBegin
ALTER TABLE moex.historical_rates DROP COLUMN short_name;
ALTER TABLE moex.index_analytics DROP COLUMN short_name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE moex.historical_rates ADD short_name NVARCHAR(100) NULL;
ALTER TABLE moex.index_analytics ADD short_name NVARCHAR(100) NULL;
-- +goose StatementEnd
