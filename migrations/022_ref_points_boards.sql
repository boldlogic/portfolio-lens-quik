--+goose Up
CREATE TABLE
    ref.trade_points (
        point_id tinyint IDENTITY (1, 1),
        code nvarchar (15) NOT NULL,
        name nvarchar (60) NOT NULL,
        CONSTRAINT PK_ref_trade_points PRIMARY KEY CLUSTERED (point_id),
    );

CREATE UNIQUE NONCLUSTERED INDEX UQ_trade_points_code ON ref.trade_points (code);

---
SET
    IDENTITY_INSERT ref.trade_points ON;

INSERT into
    ref.trade_points (point_id, code, name)
select
    src.point_id,
    src.code,
    src.name
from
    (
        VALUES
            (1, 'MOEX', N'Московская биржа'),
            (2, 'SPB', N'Санкт-петербургская биржа')
    ) AS src (point_id, code, name)
SET
    IDENTITY_INSERT ref.trade_points OFF;

---
CREATE TABLE
    ref.boards (
        board_id tinyint IDENTITY (1, 1),
        code nvarchar (12) NOT NULL,
        name nvarchar (60) NOT NULL,
        trade_point_id tinyint NULL,
        is_traded bit NOT NULL DEFAULT 0,
        CONSTRAINT PK_ref_boards PRIMARY KEY CLUSTERED (board_id),
        CONSTRAINT FK_ref_boards_trade_point_id FOREIGN KEY (trade_point_id) REFERENCES ref.trade_points (point_id),
    );

CREATE UNIQUE NONCLUSTERED INDEX UQ_boards_code ON ref.boards (code);

WITH
    src AS (
        SELECT
            s.code,
            s.name,
            s.trade_point_id,
            s.is_traded
        FROM
            (
                VALUES
                    ('TQBR', N'МБ ФР: Т+: Акции', 1, 1),
                    ('TQTF', N'МБ ФР: Т+: ETF', 1, 1),
                    ('CETS', N'МБ Валюта: ЕТС', 1, 1),
                    ('TQCB', N'МБ ФР: Т+: Корпоративные облигации', 1, 1),
                    ('CETS_MTL', N'МБ Валюта: ЕТС (Металлы)', 1, 1)
            ) AS s (code, name, trade_point_id, is_traded)
    ) MERGE INTO ref.boards AS tgt USING src ON tgt.code=src.code WHEN MATCHED
    AND (
        tgt.name<>src.name
        or tgt.trade_point_id<>src.trade_point_id
        or tgt.is_traded<>src.is_traded
    ) THEN
UPDATE
SET
    tgt.name=src.name,
    tgt.trade_point_id=src.trade_point_id,
    tgt.is_traded=src.is_traded 
WHEN NOT MATCHED BY TARGET 
THEN 
INSERT (code, name, trade_point_id, is_traded)
VALUES
    (
        src.code,
        src.name,
        src.trade_point_id,
        src.is_traded
    );

---
CREATE TABLE
    ref.instrument_types (
        type_id tinyint IDENTITY (1, 1),
        title nvarchar (150) NOT NULL,
        CONSTRAINT PK_ref_instrument_types PRIMARY KEY CLUSTERED (type_id),
        CONSTRAINT UQ_ref_instrument_types_title UNIQUE (title)
    );

SET
    IDENTITY_INSERT ref.instrument_types ON;

INSERT into
    ref.instrument_types (type_id, title)
select
    src.type_id,
    src.title
from
    (
        VALUES
            (1, 'Акции'),
            (2, 'Облигации'),
            (3, 'Валюта')
    ) AS src (type_id, title)
SET
    IDENTITY_INSERT ref.instrument_types OFF;