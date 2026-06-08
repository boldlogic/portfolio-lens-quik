-- +goose Up

CREATE SCHEMA ref;

CREATE TABLE ref.external_systems (
    ext_system_id TINYINT     NOT NULL IDENTITY(1, 1),
    ext_system    VARCHAR(50) NOT NULL,
    CONSTRAINT PK_ref_external_systems PRIMARY KEY CLUSTERED (ext_system_id)
);

CREATE UNIQUE NONCLUSTERED INDEX NCLU_ref_ext_system
    ON ref.external_systems (ext_system);

CREATE TABLE ref.currencies (
    iso_code      SMALLINT          NOT NULL,
    iso_char_code CHAR(3)           NOT NULL,
    currency_name NVARCHAR(100)     NULL,
    lat_name      NVARCHAR(100)     NULL,
    minor_units   INT               NULL,
    created_at    DATETIMEOFFSET(7) NOT NULL DEFAULT SYSDATETIMEOFFSET(),
    updated_at    DATETIMEOFFSET(7) NOT NULL DEFAULT SYSDATETIMEOFFSET(),
    ext_system_id TINYINT           NULL,
    CONSTRAINT PK_ref_currencies PRIMARY KEY CLUSTERED (iso_code),
    CONSTRAINT FK_ref_currencies_ext_system FOREIGN KEY (ext_system_id)
        REFERENCES ref.external_systems (ext_system_id)
);

CREATE TABLE ref.external_codes (
    external_code_id INT            NOT NULL IDENTITY(1, 1),
    ext_system_id    TINYINT        NOT NULL,
    ext_code         NVARCHAR(100)  NOT NULL,
    ext_code_type_id TINYINT        NOT NULL,
    internal_id      BIGINT         NOT NULL,
    CONSTRAINT PK_ref_external_codes PRIMARY KEY CLUSTERED (external_code_id),
    CONSTRAINT FK_ref_external_codes_ext_system FOREIGN KEY (ext_system_id)
        REFERENCES ref.external_systems (ext_system_id)
);

CREATE UNIQUE NONCLUSTERED INDEX NCLU_ref_external_codes_code_type_system
    ON ref.external_codes (ext_code, ext_code_type_id, ext_system_id)
    INCLUDE (internal_id);

CREATE TABLE ref.firms (
    firm_id TINYINT      NOT NULL IDENTITY(1, 1),
    code    VARCHAR(12)  NOT NULL,
    name    VARCHAR(128) NULL,
    CONSTRAINT PK_ref_firms PRIMARY KEY CLUSTERED (firm_id)
);

CREATE UNIQUE NONCLUSTERED INDEX UQ_ref_firms_code
    ON ref.firms (code);
