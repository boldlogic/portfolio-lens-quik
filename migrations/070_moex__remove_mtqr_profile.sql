-- +goose Up
-- +goose StatementBegin
-- Full-board quotes (MTQR) не поддерживаются: quotes worker требует tickers в watchlist.
DELETE FROM moex.sync_board_profiles WHERE board_id = N'MTQR';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
IF NOT EXISTS (SELECT 1 FROM moex.sync_board_profiles WHERE board_id = N'MTQR')
    INSERT INTO moex.sync_board_profiles (board_id, engine, market, quote_style, enabled)
    VALUES (N'MTQR', N'otc', N'shares', NULL, 1);
-- +goose StatementEnd
