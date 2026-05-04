SET IDENTITY_INSERT dbo.external_systems ON;

INSERT INTO dbo.external_systems (ext_system_id, ext_system)
SELECT t.ext_system_id, t.ext_system
FROM (VALUES
    (1, 'QUIK')
) AS t(ext_system_id, ext_system)
WHERE NOT EXISTS (SELECT 1 FROM dbo.external_systems e WHERE e.ext_system = t.ext_system);

SET IDENTITY_INSERT dbo.external_systems OFF;
GO
