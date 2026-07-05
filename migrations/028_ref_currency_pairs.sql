-- +goose Up

IF NOT EXISTS (
    SELECT 1
    FROM ref.instrument_types
    WHERE title = N'Валютные пары'
)
INSERT INTO ref.instrument_types (title)
VALUES (N'Валютные пары');

-- +goose StatementBegin
DECLARE @currency_pair_type_id TINYINT = (
    SELECT type_id
    FROM ref.instrument_types
    WHERE title = N'Валютные пары'
);

IF @currency_pair_type_id IS NULL
BEGIN
    THROW 51002, N'Тип инструмента Валютные пары не найден', 1;
END;

INSERT INTO ref.instrument_type_boards (type_id, board_id)
SELECT @currency_pair_type_id, b.board_id
FROM ref.boards b
WHERE b.code IN (N'CETS', N'CETS_MTL')
    AND NOT EXISTS (
        SELECT 1
        FROM ref.instrument_type_boards itb
        WHERE itb.board_id = b.board_id
    );
-- +goose StatementEnd
