IF OBJECT_ID (N'dbo.firms', N'U') IS NULL BEGIN
CREATE TABLE
    dbo.firms (
        firm_id tinyint IDENTITY (1, 1),
        code varchar(12) NOT NULL,
        name varchar(128) CONSTRAINT PK_quik_firms PRIMARY KEY CLUSTERED (firm_id),
    );

END 
GO 
IF OBJECT_ID (N'dbo.firms', N'U') IS NOT NULL 
BEGIN
DROP INDEX IF EXISTS UQ_firms_code ON dbo.firms;

CREATE UNIQUE NONCLUSTERED INDEX UQ_firms_code ON dbo.firms (code);

END 

GO