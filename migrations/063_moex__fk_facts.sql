-- +goose Up
-- +goose StatementBegin
DELETE FROM moex.historical_rates
WHERE sec_id NOT IN (SELECT sec_id FROM moex.securities);

DELETE FROM moex.historical_rates
WHERE board_id NOT IN (SELECT board_id FROM moex.boards);

DELETE FROM moex.index_analytics
WHERE sec_id IS NOT NULL
  AND sec_id NOT IN (SELECT sec_id FROM moex.securities);

ALTER TABLE moex.historical_rates
    ADD CONSTRAINT FK_historical_rates_securities
    FOREIGN KEY (sec_id) REFERENCES moex.securities(sec_id);

ALTER TABLE moex.historical_rates
    ADD CONSTRAINT FK_historical_rates_boards
    FOREIGN KEY (board_id) REFERENCES moex.boards(board_id);

ALTER TABLE moex.index_analytics
    ADD CONSTRAINT FK_index_analytics_securities
    FOREIGN KEY (sec_id) REFERENCES moex.securities(sec_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE moex.index_analytics DROP CONSTRAINT FK_index_analytics_securities;
ALTER TABLE moex.historical_rates DROP CONSTRAINT FK_historical_rates_boards;
ALTER TABLE moex.historical_rates DROP CONSTRAINT FK_historical_rates_securities;
-- +goose StatementEnd
