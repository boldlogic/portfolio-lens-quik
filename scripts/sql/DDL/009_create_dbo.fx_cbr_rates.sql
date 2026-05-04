IF OBJECT_ID(N'dbo.fx_cbr_rates', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.fx_cbr_rates (
        date                date              NOT NULL,
        quote_iso_code      SMALLINT               NOT NULL,
        base_iso_code       SMALLINT               NOT NULL,

        rate_quote_per_base      DECIMAL(18,8)     NULL,
        rate_base_per_quote      DECIMAL(18,8)     NULL,
        created_at          DATETIMEOFFSET(7) NOT NULL DEFAULT SYSDATETIMEOFFSET(),
        updated_at          DATETIMEOFFSET(7) NOT NULL DEFAULT SYSDATETIMEOFFSET(),
        ext_system_id tinyint NULL,
        CONSTRAINT PK_fx_cbr_rates PRIMARY KEY CLUSTERED (date, quote_iso_code, base_iso_code),
        CONSTRAINT FK_fx_cbr_rates_quote_iso FOREIGN KEY (quote_iso_code) REFERENCES dbo.currencies (iso_code),
        CONSTRAINT FK_fx_cbr_rates_base_iso FOREIGN KEY (quote_iso_code) REFERENCES dbo.currencies (iso_code),
        CONSTRAINT FK_fx_cbr_rates_ext_system FOREIGN KEY (ext_system_id) REFERENCES dbo.external_systems (ext_system_id)
    );
END
GO
