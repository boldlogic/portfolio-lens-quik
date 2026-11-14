USE [master];

DECLARE @app_db SYSNAME = N'portfolio_lens_quik';
DECLARE @login_name SYSNAME = N't_invest_portfolio_app';
DECLARE @user_name SYSNAME = N'u_t_invest_portfolio_app';
DECLARE @role_name SYSNAME = N't_invest_portfolio_writer';
DECLARE @password NVARCHAR(256) = N'CHANGE_ME';

IF DB_ID(@app_db) IS NULL
    THROW 51000, N'База данных portfolio_lens_quik не существует', 1;

IF @password = N'CHANGE_ME'
    THROW 51001, N'Задайте пароль в локальной копии bootstrap-скрипта', 1;

IF NOT EXISTS (SELECT 1 FROM sys.sql_logins WHERE name = @login_name)
BEGIN
    DECLARE @create_login_sql NVARCHAR(MAX) =
        N'CREATE LOGIN ' + QUOTENAME(@login_name) +
        N' WITH PASSWORD = ' + QUOTENAME(@password, '''') +
        N', CHECK_POLICY = ON, CHECK_EXPIRATION = OFF, DEFAULT_DATABASE = ' + QUOTENAME(@app_db) + N';';
    EXEC (@create_login_sql);
END;

DECLARE @create_database_principals_sql NVARCHAR(MAX) =
N'USE ' + QUOTENAME(@app_db) + N';
IF NOT EXISTS (SELECT 1 FROM sys.database_principals WHERE name = N' + QUOTENAME(@user_name, '''') + N')
    CREATE USER ' + QUOTENAME(@user_name) + N' FOR LOGIN ' + QUOTENAME(@login_name) + N';
IF NOT EXISTS (SELECT 1 FROM sys.database_principals WHERE name = N' + QUOTENAME(@role_name, '''') + N' AND type = ''R'')
    CREATE ROLE ' + QUOTENAME(@role_name) + N' AUTHORIZATION dbo;
IF NOT EXISTS (
    SELECT 1
    FROM sys.database_role_members drm
    JOIN sys.database_principals r ON r.principal_id = drm.role_principal_id
    JOIN sys.database_principals m ON m.principal_id = drm.member_principal_id
    WHERE r.name = N' + QUOTENAME(@role_name, '''') + N'
      AND m.name = N' + QUOTENAME(@user_name, '''') + N'
)
    ALTER ROLE ' + QUOTENAME(@role_name) + N' ADD MEMBER ' + QUOTENAME(@user_name) + N';';
EXEC (@create_database_principals_sql);
