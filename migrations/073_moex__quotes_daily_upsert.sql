-- +goose Up
-- +goose StatementBegin
ALTER TABLE moex.quotes ADD quote_date DATE NULL;
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE moex.quotes
SET quote_date = CAST(snapshot_at AT TIME ZONE 'UTC' AT TIME ZONE 'Russian Standard Time' AS DATE)
WHERE quote_date IS NULL;
-- +goose StatementEnd

-- +goose StatementBegin
;WITH ranked AS (
    SELECT id,
           ROW_NUMBER() OVER (
               PARTITION BY sec_id, board_id, quote_date
               ORDER BY snapshot_at DESC
           ) AS rn
    FROM moex.quotes
    WHERE quote_date IS NOT NULL
)
DELETE FROM moex.quotes
WHERE id IN (SELECT id FROM ranked WHERE rn > 1);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE moex.quotes ALTER COLUMN quote_date DATE NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX UQ_quotes_sec_board_date ON moex.quotes(sec_id, board_id, quote_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS UQ_quotes_sec_board_date ON moex.quotes;
ALTER TABLE moex.quotes DROP COLUMN quote_date;
-- +goose StatementEnd
