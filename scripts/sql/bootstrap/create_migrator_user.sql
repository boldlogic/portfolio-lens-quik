USE [master];

DECLARE @migrator_db SYSNAME = N'portfolio_lens_quik';
DECLARE @migrator_login SYSNAME = N'quik_portfolio_migrator';
DECLARE @migrator_user SYSNAME = N'quik_portfolio_migrator';
DECLARE @migrator_role SYSNAME = N'quik_portfolio_admin';
DECLARE @migrator_password NVARCHAR(256) = N'CHANGE_ME';

IF DB_ID(@migrator_db) IS NULL
BEGIN
    RAISERROR(N'База данных не существует', 16, 1);
    RETURN;
END;

IF NOT EXISTS (
    SELECT 1
    FROM sys.sql_logins
    WHERE name = @migrator_login
)
BEGIN
    IF @migrator_password = N'CHANGE_ME'
    BEGIN
        RAISERROR(N'Укажите пароль в переменной @migrator_password', 16, 1);
        RETURN;
    END;

    DECLARE @create_login_sql NVARCHAR(MAX) =
        N'CREATE LOGIN ' + QUOTENAME(@migrator_login) + N' WITH PASSWORD = ' + QUOTENAME(@migrator_password, '''') + N';';
    EXEC (@create_login_sql);
END;

DECLARE @set_default_language_sql NVARCHAR(MAX) =
    N'ALTER LOGIN ' + QUOTENAME(@migrator_login) + N' WITH DEFAULT_LANGUAGE = [Russian];';
EXEC (@set_default_language_sql);

DECLARE @create_user_sql NVARCHAR(MAX) =
N'USE ' + QUOTENAME(@migrator_db) + N';
IF NOT EXISTS (
    SELECT 1
    FROM sys.database_principals
    WHERE name = N' + QUOTENAME(@migrator_user, '''') + N'
)
BEGIN
    CREATE USER ' + QUOTENAME(@migrator_user) + N' FOR LOGIN ' + QUOTENAME(@migrator_login) + N';
END;';
EXEC (@create_user_sql);

DECLARE @create_role_sql NVARCHAR(MAX) =
N'USE ' + QUOTENAME(@migrator_db) + N';
IF NOT EXISTS (
    SELECT 1
    FROM sys.database_principals
    WHERE name = N' + QUOTENAME(@migrator_role, '''') + N'
      AND type = ''R''
)
BEGIN
    CREATE ROLE ' + QUOTENAME(@migrator_role) + N';
END;';
EXEC (@create_role_sql);

DECLARE @add_member_sql NVARCHAR(MAX) =
N'USE ' + QUOTENAME(@migrator_db) + N';
IF NOT EXISTS (
    SELECT 1
    FROM sys.database_role_members drm
    JOIN sys.database_principals r ON r.principal_id = drm.role_principal_id
    JOIN sys.database_principals m ON m.principal_id = drm.member_principal_id
    WHERE r.name = N' + QUOTENAME(@migrator_role, '''') + N'
      AND m.name = N' + QUOTENAME(@migrator_user, '''') + N'
)
BEGIN
    ALTER ROLE ' + QUOTENAME(@migrator_role) + N' ADD MEMBER ' + QUOTENAME(@migrator_user) + N';
END;';
EXEC (@add_member_sql);

DECLARE @add_ddladmin_sql NVARCHAR(MAX) =
N'USE ' + QUOTENAME(@migrator_db) + N';
IF NOT EXISTS (
    SELECT 1
    FROM sys.database_role_members drm
    JOIN sys.database_principals r ON r.principal_id = drm.role_principal_id
    JOIN sys.database_principals m ON m.principal_id = drm.member_principal_id
    WHERE r.name = N''db_ddladmin''
      AND m.name = N' + QUOTENAME(@migrator_role, '''') + N'
)
BEGIN
    ALTER ROLE db_ddladmin ADD MEMBER ' + QUOTENAME(@migrator_role) + N';
END;';
EXEC (@add_ddladmin_sql);

DECLARE @add_datareader_sql NVARCHAR(MAX) =
N'USE ' + QUOTENAME(@migrator_db) + N';
IF NOT EXISTS (
    SELECT 1
    FROM sys.database_role_members drm
    JOIN sys.database_principals r ON r.principal_id = drm.role_principal_id
    JOIN sys.database_principals m ON m.principal_id = drm.member_principal_id
    WHERE r.name = N''db_datareader''
      AND m.name = N' + QUOTENAME(@migrator_role, '''') + N'
)
BEGIN
    ALTER ROLE db_datareader ADD MEMBER ' + QUOTENAME(@migrator_role) + N';
END;';
EXEC (@add_datareader_sql);

DECLARE @add_datawriter_sql NVARCHAR(MAX) =
N'USE ' + QUOTENAME(@migrator_db) + N';
IF NOT EXISTS (
    SELECT 1
    FROM sys.database_role_members drm
    JOIN sys.database_principals r ON r.principal_id = drm.role_principal_id
    JOIN sys.database_principals m ON m.principal_id = drm.member_principal_id
    WHERE r.name = N''db_datawriter''
      AND m.name = N' + QUOTENAME(@migrator_role, '''') + N'
)
BEGIN
    ALTER ROLE db_datawriter ADD MEMBER ' + QUOTENAME(@migrator_role) + N';
END;';
EXEC (@add_datawriter_sql);
