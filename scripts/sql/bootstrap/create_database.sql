USE [master];

IF DB_ID(N'portfolio_lens_quik') IS NULL
BEGIN
    CREATE DATABASE [portfolio_lens_quik];
    ALTER DATABASE [portfolio_lens_quik] COLLATE Cyrillic_General_CI_AS;
END;
