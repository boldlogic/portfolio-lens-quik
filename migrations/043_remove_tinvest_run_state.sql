-- +goose Up
-- +goose StatementBegin
IF OBJECT_ID(N'tinvest.snapshot_warning', N'U') IS NOT NULL DROP TABLE tinvest.snapshot_warning;
IF OBJECT_ID(N'tinvest.snapshot_stop_order', N'U') IS NOT NULL DROP TABLE tinvest.snapshot_stop_order;
IF OBJECT_ID(N'tinvest.snapshot_order', N'U') IS NOT NULL DROP TABLE tinvest.snapshot_order;
IF OBJECT_ID(N'tinvest.virtual_position', N'U') IS NOT NULL DROP TABLE tinvest.virtual_position;
IF OBJECT_ID(N'tinvest.holding', N'U') IS NOT NULL DROP TABLE tinvest.holding;
IF OBJECT_ID(N'tinvest.position', N'U') IS NOT NULL DROP TABLE tinvest.position;
IF OBJECT_ID(N'tinvest.money_amount', N'U') IS NOT NULL DROP TABLE tinvest.money_amount;
IF OBJECT_ID(N'tinvest.percent_reading', N'U') IS NOT NULL DROP TABLE tinvest.percent_reading;
IF OBJECT_ID(N'tinvest.snapshot_account', N'U') IS NOT NULL DROP TABLE tinvest.snapshot_account;
IF OBJECT_ID(N'tinvest.snapshot', N'U') IS NOT NULL DROP TABLE tinvest.snapshot;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
THROW 51043, N'Irreversible migration: legacy run-scoped state was removed', 1;
-- +goose StatementEnd
