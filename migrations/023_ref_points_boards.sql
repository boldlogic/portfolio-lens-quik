--+goose Up

INSERT into
    ref.boards (code, name, trade_point_id, is_traded)
select
    src.code,
    src.name,
    src.trade_point_id,
    src.is_traded
from
    (
        VALUES
            ('CROSSRATE', N'Кросс-курсы валют', null, 0),
            ('SMAL', N'МБ ФР: Т+: Неполные лоты', 1, 1),
            ('SPBRU', N'SPB: Российские  Акции', 2, 1)
    ) 
AS src (code, name, trade_point_id, is_traded)

---
CREATE TABLE
    ref.instrument_type_boards (
        type_id tinyint,
        board_id tinyint,
        CONSTRAINT PK_ref_instrument_type_boards PRIMARY KEY CLUSTERED (type_id, board_id),
        CONSTRAINT UQ_ref_instrument_type_boards_board_id UNIQUE (board_id),
        CONSTRAINT FK_ref_instrument_type_boards_type_id FOREIGN KEY (type_id) REFERENCES ref.instrument_types (type_id),
        CONSTRAINT FK_ref_instrument_type_boards_board_id FOREIGN KEY (board_id) REFERENCES ref.boards (board_id)
    );

with src as(
select v.board, v.type_id, b.board_id
from
    (
        VALUES
            ('TQBR', 1),
			('TQTF', 1),
			('TQCB', 2),
			('CROSSRATE', 3),
            ('SMAL', 1),
            ('SPBRU', 1)
    ) 
AS v (board, type_id)
join ref.boards b on b.code=v.board
)
insert into ref.instrument_type_boards (type_id, board_id)
select src.type_id, src.board_id from src;

