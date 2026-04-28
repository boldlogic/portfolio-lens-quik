IF OBJECT_ID (N'quik.firms', N'U') IS NULL BEGIN
CREATE TABLE
    quik.firms (
        firm_id tinyint IDENTITY (1, 1),
        code varchar(12) NOT NULL,
        name varchar(128) CONSTRAINT PK_quik_firms PRIMARY KEY CLUSTERED (firm_id),
    );

END 
GO 
IF OBJECT_ID (N'quik.firms', N'U') IS NOT NULL 
BEGIN
DROP INDEX IF EXISTS UQ_firms_code ON quik.firms;

CREATE UNIQUE NONCLUSTERED INDEX UQ_firms_code ON quik.firms (code);

END 

GO