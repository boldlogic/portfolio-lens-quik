IF NOT EXISTS (
    SELECT 1
    FROM sys.schemas
    WHERE name = N'api'
)
BEGIN
    EXEC(N'CREATE SCHEMA [api] AUTHORIZATION [quik_portfolio_app]');
END;
