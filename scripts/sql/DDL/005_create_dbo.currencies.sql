IF OBJECT_ID(N'quik.currencies', N'U') IS NULL
BEGIN
    CREATE TABLE quik.currencies (
        iso_code       SMALLINT          NOT NULL,
        iso_char_code  CHAR(3)       NOT NULL,
        currency_name  NVARCHAR(100)     NULL,
        lat_name       NVARCHAR(100)     NULL,
        minor_units                      INT               NULL,
        created_at                       DATETIMEOFFSET(7) NOT NULL DEFAULT SYSDATETIMEOFFSET(),
        updated_at                       DATETIMEOFFSET(7) NOT NULL DEFAULT SYSDATETIMEOFFSET(),
        CONSTRAINT PK_currencies PRIMARY KEY CLUSTERED (iso_code)
            );
END
GO