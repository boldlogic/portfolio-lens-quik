-- +goose Up

-- Текущая позиция уникальна в пределах серверного календарного дня, а не запуска.
-- При схлопывании старых повторных запусков сохраняется последняя строка каждого дня/счёта/инструмента.
-- Для ранее накопленных строк день пересчитывается к текущему локальному offset SQL Server.
-- +goose StatementBegin
DECLARE @fk sysname;
SELECT TOP (1) @fk = fk.name
FROM sys.foreign_keys fk
JOIN sys.foreign_key_columns fkc ON fkc.constraint_object_id = fk.object_id
JOIN sys.columns c ON c.object_id = fkc.parent_object_id AND c.column_id = fkc.parent_column_id
WHERE fk.parent_object_id = OBJECT_ID(N'tinvest.position') AND c.name = N'snapshot_id';
IF @fk IS NOT NULL BEGIN DECLARE @sql nvarchar(max) = N'ALTER TABLE tinvest.position DROP CONSTRAINT ' + QUOTENAME(@fk); EXEC sys.sp_executesql @sql; END;

DECLARE @pk sysname;
SELECT @pk = kc.name FROM sys.key_constraints kc
WHERE kc.parent_object_id = OBJECT_ID(N'tinvest.position') AND kc.type = 'PK';
IF @pk IS NOT NULL BEGIN DECLARE @sql2 nvarchar(max) = N'ALTER TABLE tinvest.position DROP CONSTRAINT ' + QUOTENAME(@pk); EXEC sys.sp_executesql @sql2; END;
-- +goose StatementEnd

ALTER TABLE tinvest.position ADD load_date DATE NULL;

UPDATE p
SET load_date = CAST(SWITCHOFFSET(s.checked_at, DATEPART(TZOFFSET, SYSDATETIMEOFFSET())) AS DATE)
FROM tinvest.position p
JOIN tinvest.snapshot s ON s.id = p.snapshot_id;

-- +goose StatementBegin
;WITH ranked AS (
    SELECT snapshot_id, account_id, instrument_uid,
           ROW_NUMBER() OVER (PARTITION BY load_date, account_id, instrument_uid ORDER BY snapshot_id DESC) AS rn
    FROM tinvest.position
)
DELETE FROM ranked WHERE rn > 1;
-- +goose StatementEnd

ALTER TABLE tinvest.position ALTER COLUMN load_date DATE NOT NULL;
ALTER TABLE tinvest.position DROP COLUMN snapshot_id;
ALTER TABLE tinvest.position ADD CONSTRAINT PK_tinvest_position PRIMARY KEY CLUSTERED (load_date, account_id, instrument_uid);

-- +goose Down
-- Схлопнутые повторы внутри одного дня нельзя восстановить как исходные запуски.
-- Запрет отката не позволяет создать ложную историю и потерять исходную семантику данных.
-- +goose StatementBegin
THROW 51041, N'Откат невозможен: повторные внутридневные записи позиций были объединены по дневному ключу', 1;
-- +goose StatementEnd
