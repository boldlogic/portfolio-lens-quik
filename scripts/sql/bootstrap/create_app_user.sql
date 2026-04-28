USE [master];

DECLARE @app_db SYSNAME = N'portfolio_lens_quik';
DECLARE @app_login SYSNAME = N'quik_portfolio_app';
DECLARE @app_user SYSNAME = N'quik_portfolio_app';
DECLARE @app_role SYSNAME = N'quik_portfolio_rw';
DECLARE @app_password NVARCHAR(256) = N'CHANGE_ME';

IF @app_password = N'CHANGE_ME'
BEGIN
    RAISERROR(N'Укажите пароль в переменной @app_password', 16, 1);
    RETURN;
END;

IF DB_ID(@app_db) IS NULL
BEGIN
    RAISERROR(N'База данных не существует', 16, 1);
    RETURN;
END;

IF NOT EXISTS (
    SELECT 1
    FROM sys.sql_logins
    WHERE name = @app_login
)
BEGIN
    DECLARE @create_login_sql NVARCHAR(MAX) =
        N'CREATE LOGIN ' + QUOTENAME(@app_login) + N' WITH PASSWORD = ' + QUOTENAME(@app_password, '''') + N';';
    EXEC (@create_login_sql);
END;

DECLARE @set_default_language_sql NVARCHAR(MAX) =
    N'ALTER LOGIN ' + QUOTENAME(@app_login) + N' WITH DEFAULT_LANGUAGE = [Russian];';
EXEC (@set_default_language_sql);

DECLARE @create_user_sql NVARCHAR(MAX) =
N'USE ' + QUOTENAME(@app_db) + N';
IF NOT EXISTS (
    SELECT 1
    FROM sys.database_principals
    WHERE name = N' + QUOTENAME(@app_user, '''') + N'
)
BEGIN
    CREATE USER ' + QUOTENAME(@app_user) + N' FOR LOGIN ' + QUOTENAME(@app_login) + N';
END;';
EXEC (@create_user_sql);

DECLARE @create_role_sql NVARCHAR(MAX) =
N'USE ' + QUOTENAME(@app_db) + N';
IF NOT EXISTS (
    SELECT 1
    FROM sys.database_principals
    WHERE name = N' + QUOTENAME(@app_role, '''') + N'
      AND type = ''R''
)
BEGIN
    CREATE ROLE ' + QUOTENAME(@app_role) + N';
END;';
EXEC (@create_role_sql);

DECLARE @add_member_sql NVARCHAR(MAX) =
N'USE ' + QUOTENAME(@app_db) + N';
IF NOT EXISTS (
    SELECT 1
    FROM sys.database_role_members drm
    JOIN sys.database_principals r ON r.principal_id = drm.role_principal_id
    JOIN sys.database_principals m ON m.principal_id = drm.member_principal_id
    WHERE r.name = N' + QUOTENAME(@app_role, '''') + N'
      AND m.name = N' + QUOTENAME(@app_user, '''') + N'
)
BEGIN
    ALTER ROLE ' + QUOTENAME(@app_role) + N' ADD MEMBER ' + QUOTENAME(@app_user) + N';
END;';
EXEC (@add_member_sql);

DECLARE @grant_sql NVARCHAR(MAX) =
N'USE ' + QUOTENAME(@app_db) + N';
GRANT SELECT, INSERT, UPDATE, DELETE ON SCHEMA::[quik] TO ' + QUOTENAME(@app_user) + N';
DENY ALTER ON SCHEMA::[quik] TO ' + QUOTENAME(@app_user) + N';
DENY CREATE TABLE TO ' + QUOTENAME(@app_user) + N';
DENY CREATE VIEW TO ' + QUOTENAME(@app_user) + N';
DENY CREATE PROCEDURE TO ' + QUOTENAME(@app_user) + N';
DENY CREATE FUNCTION TO ' + QUOTENAME(@app_user) + N';
DENY CREATE TYPE TO ' + QUOTENAME(@app_user) + N';
DENY CREATE SCHEMA TO ' + QUOTENAME(@app_user) + N';';
EXEC (@grant_sql);
