IF OBJECT_ID(N'dbo.currencies', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.currencies (
        iso_code       SMALLINT          NOT NULL,
        iso_char_code  CHAR(3)       NOT NULL,
        currency_name  NVARCHAR(100)     NULL,
        lat_name       NVARCHAR(100)     NULL,
        minor_units                      INT               NULL,
        created_at                       DATETIMEOFFSET(7) NOT NULL DEFAULT SYSDATETIMEOFFSET(),
        updated_at                       DATETIMEOFFSET(7) NOT NULL DEFAULT SYSDATETIMEOFFSET(),
        ext_system_id tinyint NULL,
        CONSTRAINT PK_currencies PRIMARY KEY CLUSTERED (iso_code),
        CONSTRAINT FK_currencies_ext_system FOREIGN KEY (ext_system_id) REFERENCES dbo.external_systems (ext_system_id)
            );
END
GO