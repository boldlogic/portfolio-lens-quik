-- +goose Up

-- +goose StatementBegin
DECLARE @description_default SYSNAME = (
    SELECT name FROM sys.default_constraints
    WHERE parent_object_id = OBJECT_ID(N'tinvest.brand')
      AND parent_column_id = COLUMNPROPERTY(OBJECT_ID(N'tinvest.brand'), N'description', 'ColumnId')
);
DECLARE @info_default SYSNAME = (
    SELECT name FROM sys.default_constraints
    WHERE parent_object_id = OBJECT_ID(N'tinvest.brand')
      AND parent_column_id = COLUMNPROPERTY(OBJECT_ID(N'tinvest.brand'), N'info', 'ColumnId')
);
DECLARE @drop_default_sql NVARCHAR(4000);
IF @description_default IS NOT NULL
BEGIN
    SET @drop_default_sql = N'ALTER TABLE tinvest.brand DROP CONSTRAINT ' + QUOTENAME(@description_default);
    EXEC sys.sp_executesql @drop_default_sql;
END;
IF @info_default IS NOT NULL
BEGIN
    SET @drop_default_sql = N'ALTER TABLE tinvest.brand DROP CONSTRAINT ' + QUOTENAME(@info_default);
    EXEC sys.sp_executesql @drop_default_sql;
END;
-- +goose StatementEnd

ALTER TABLE tinvest.brand ALTER COLUMN description NVARCHAR(MAX) NOT NULL;
ALTER TABLE tinvest.brand ALTER COLUMN info NVARCHAR(MAX) NOT NULL;
ALTER TABLE tinvest.brand ADD CONSTRAINT DF_tinvest_brand_description DEFAULT N'' FOR description;
ALTER TABLE tinvest.brand ADD CONSTRAINT DF_tinvest_brand_info DEFAULT N'' FOR info;

-- +goose Down

-- +goose StatementBegin
IF EXISTS (SELECT 1 FROM tinvest.brand WHERE LEN(description) > 1000 OR LEN(info) > 2000)
    THROW 51040, N'Есть описания брендов длиннее исходного размера; сначала выгрузите или сократите их', 1;
-- +goose StatementEnd

ALTER TABLE tinvest.brand DROP CONSTRAINT DF_tinvest_brand_description;
ALTER TABLE tinvest.brand DROP CONSTRAINT DF_tinvest_brand_info;
ALTER TABLE tinvest.brand ALTER COLUMN description NVARCHAR(1000) NOT NULL;
ALTER TABLE tinvest.brand ALTER COLUMN info NVARCHAR(2000) NOT NULL;
ALTER TABLE tinvest.brand ADD CONSTRAINT DF_tinvest_brand_description DEFAULT N'' FOR description;
ALTER TABLE tinvest.brand ADD CONSTRAINT DF_tinvest_brand_info DEFAULT N'' FOR info;
