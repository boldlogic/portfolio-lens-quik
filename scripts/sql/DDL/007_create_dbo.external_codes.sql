IF OBJECT_ID(N'dbo.external_codes', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.external_codes (
        external_code_id INT           NOT NULL IDENTITY(1, 1),
        ext_system_id    TINYINT       NOT NULL,
        ext_code         NVARCHAR(100) NOT NULL,
        ext_code_type_id tinyint           NOT NULL,  -- 1=currency, 2=instrument
        internal_id      BIGINT           NOT NULL,
        CONSTRAINT PK_external_codes PRIMARY KEY CLUSTERED (external_code_id),
        CONSTRAINT FK_external_codes_ext_system FOREIGN KEY (ext_system_id) REFERENCES dbo.external_systems (ext_system_id)
    );
END
GO

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'NCLU_external_codes_code_type_system' AND object_id = OBJECT_ID(N'dbo.external_codes'))
BEGIN
    CREATE UNIQUE NONCLUSTERED INDEX NCLU_external_codes_code_type_system
        ON dbo.external_codes (ext_code, ext_code_type_id, ext_system_id)
        INCLUDE (internal_id);
END
GO
