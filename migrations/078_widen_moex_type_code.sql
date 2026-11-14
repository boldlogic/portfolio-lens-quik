-- +goose Up
-- +goose StatementBegin

ALTER TABLE moex.securities DROP CONSTRAINT FK_securities_security_types;

DECLARE @pk_name SYSNAME = (
    SELECT kc.name
    FROM sys.key_constraints kc
    WHERE kc.parent_object_id = OBJECT_ID(N'moex.security_types') AND kc.type = 'PK'
);
IF @pk_name IS NOT NULL
BEGIN
    DECLARE @drop_pk_sql NVARCHAR(300) = N'ALTER TABLE moex.security_types DROP CONSTRAINT ' + QUOTENAME(@pk_name);
    EXEC sys.sp_executesql @drop_pk_sql;
END;

ALTER TABLE moex.security_types ALTER COLUMN type_code NVARCHAR(50) NOT NULL;
ALTER TABLE moex.securities ALTER COLUMN type_code NVARCHAR(50) NULL;

ALTER TABLE moex.security_types
    ADD CONSTRAINT PK_moex_security_types PRIMARY KEY CLUSTERED (type_code);

ALTER TABLE moex.securities
    ADD CONSTRAINT FK_securities_security_types FOREIGN KEY (type_code)
        REFERENCES moex.security_types (type_code);
-- +goose StatementEnd

-- +goose Down

THROW 51000, N'Новые коды типов могут превышать старую ширину столбца; миграция не откатывается автоматически.', 1;
