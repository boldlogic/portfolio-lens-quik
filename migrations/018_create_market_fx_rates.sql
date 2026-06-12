-- +goose Up
CREATE TABLE market.fx_rates (
    date               DATE              NOT NULL,
    quote_iso_code     SMALLINT          NOT NULL,
    base_iso_code      SMALLINT          NOT NULL,
    rate DECIMAL(18, 8)   NULL,
    cb_rate DECIMAL(18, 8)   NULL,
    created_at         DATETIMEOFFSET(7) NOT NULL DEFAULT SYSDATETIMEOFFSET(),
    updated_at         DATETIMEOFFSET(7) NOT NULL DEFAULT SYSDATETIMEOFFSET(),
    ext_system_id      TINYINT           NULL,
    CONSTRAINT PK_market_fx_rates PRIMARY KEY CLUSTERED (date, quote_iso_code, base_iso_code),
    CONSTRAINT FK_market_fx_rates_quote_iso FOREIGN KEY (quote_iso_code)
        REFERENCES ref.currencies (iso_code),
    CONSTRAINT FK_market_fx_rates_base_iso FOREIGN KEY (base_iso_code)
        REFERENCES ref.currencies (iso_code),
    CONSTRAINT FK_market_fx_rates_ext_system FOREIGN KEY (ext_system_id)
        REFERENCES ref.external_systems (ext_system_id)
);
 
