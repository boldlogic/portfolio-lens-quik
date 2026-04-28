IF NOT EXISTS (
    SELECT 1
    FROM sys.schemas
    WHERE name = N'quik'
)
BEGIN
    EXEC(N'CREATE SCHEMA [quik]');
END;
