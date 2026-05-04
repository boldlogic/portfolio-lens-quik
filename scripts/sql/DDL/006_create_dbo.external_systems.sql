IF OBJECT_ID(N'dbo.external_systems', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.external_systems (
        ext_system_id   TINYINT       NOT NULL IDENTITY(1, 1),
        ext_system      VARCHAR(50)   NOT NULL,
        CONSTRAINT PK_external_systems PRIMARY KEY CLUSTERED (ext_system_id)
    );
END
GO

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'NCLU_ext_system' AND object_id = OBJECT_ID(N'dbo.external_systems'))
BEGIN
    CREATE UNIQUE NONCLUSTERED INDEX NCLU_ext_system ON dbo.external_systems (ext_system);
END
GO
