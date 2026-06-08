USE [master];

DECLARE @app_db SYSNAME = N'portfolio_lens_quik';
DECLARE @login_name SYSNAME = N'quik_portfolio_writer';
DECLARE @user_name SYSNAME = N'u_quik_portfolio_writer';
DECLARE @role_name SYSNAME = N'quik_portfolio_writer';
DECLARE @password NVARCHAR(256) = N'CHANGE_ME';

IF DB_ID(@app_db) IS NULL
BEGIN
    RAISERROR(N'База данных не существует', 16, 1);
    RETURN;
END;

IF NOT EXISTS (
    SELECT 1
    FROM sys.sql_logins
    WHERE name = @login_name
)
BEGIN
    IF @password = N'CHANGE_ME'
    BEGIN
        RAISERROR(N'Укажите пароль в переменной @password', 16, 1);
        RETURN;
    END;

    DECLARE @create_login_sql NVARCHAR(MAX) =
        N'CREATE LOGIN ' + QUOTENAME(@login_name) + N' WITH PASSWORD = ' + QUOTENAME(@password, '''') + N';';
    EXEC (@create_login_sql);
END;

DECLARE @set_default_language_sql NVARCHAR(MAX) =
    N'ALTER LOGIN ' + QUOTENAME(@login_name) + N' WITH DEFAULT_LANGUAGE = [Russian];';
EXEC (@set_default_language_sql);

DECLARE @create_user_sql NVARCHAR(MAX) =
N'USE ' + QUOTENAME(@app_db) + N';
IF NOT EXISTS (
    SELECT 1
    FROM sys.database_principals
    WHERE name = N' + QUOTENAME(@user_name, '''') + N'
)
BEGIN
    CREATE USER ' + QUOTENAME(@user_name) + N' FOR LOGIN ' + QUOTENAME(@login_name) + N';
END;

IF NOT EXISTS (
    SELECT 1
    FROM sys.database_principals
    WHERE name = N' + QUOTENAME(@role_name, '''') + N'
      AND type = ''R''
)
BEGIN
    CREATE ROLE ' + QUOTENAME(@role_name) + N';
END;

IF NOT EXISTS (
    SELECT 1
    FROM sys.database_role_members drm
    JOIN sys.database_principals r ON r.principal_id = drm.role_principal_id
    JOIN sys.database_principals m ON m.principal_id = drm.member_principal_id
    WHERE r.name = N' + QUOTENAME(@role_name, '''') + N'
      AND m.name = N' + QUOTENAME(@user_name, '''') + N'
)
BEGIN
    ALTER ROLE ' + QUOTENAME(@role_name) + N' ADD MEMBER ' + QUOTENAME(@user_name) + N';
END;';
EXEC (@create_user_sql);
