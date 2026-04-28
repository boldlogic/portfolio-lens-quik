IF NOT EXISTS (
    SELECT
        *
    FROM
        sys.tables t
        INNER JOIN sys.schemas s ON t.schema_id=s.schema_id
    WHERE
        s.name=N'quik'
        AND t.name=N'security_limits_otc'
) BEGIN
CREATE TABLE
    quik.security_limits_otc (
        load_date date NOT NULL DEFAULT (getdate ()),
        client_code varchar(12) NOT NULL,
        ticker varchar(12) NOT NULL,
        trade_account varchar(12) NOT NULL DEFAULT 'OTC',
        settle_code varchar(5) NOT NULL DEFAULT 'Tx',
        firm_code varchar(12) NOT NULL,
        firm_name varchar(128) NULL,
        balance DECIMAL(19,4) NULL,
        acquisition_ccy varchar(3) NULL,
        isin varchar(12) NULL,
        source_date date NOT NULL DEFAULT (getdate()),
        ts timestamp NOT NULL,
        
        CONSTRAINT PK_quik_security_limits_otc PRIMARY KEY CLUSTERED (
            load_date ASC,
            client_code ASC,
            ticker ASC,
            trade_account ASC,
            settle_code ASC,
            firm_code ASC
        ),
    );

END 
GO
