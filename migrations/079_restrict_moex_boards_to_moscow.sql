-- +goose Up

DELETE m
FROM moex.common_board_map m
JOIN ref.boards b ON b.board_id = m.board_id
WHERE b.trade_point_id IS NULL OR b.trade_point_id <> 1;

-- +goose Down

THROW 51000, N'Откат может восстановить связь с торговой точкой, не относящейся к Московской бирже.', 1;
