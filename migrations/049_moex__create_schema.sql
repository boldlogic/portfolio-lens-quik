-- +goose Up
-- +goose StatementBegin
IF NOT EXISTS (SELECT 1 FROM sys.schemas WHERE name = N'moex')
    EXEC(N'CREATE SCHEMA moex');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
IF EXISTS (SELECT 1 FROM sys.schemas WHERE name = N'moex')
    EXEC(N'DROP SCHEMA moex');
-- +goose StatementEnd
