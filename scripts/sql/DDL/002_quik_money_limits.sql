IF NOT EXISTS (
    SELECT
        *
    FROM
        sys.tables t
        INNER JOIN sys.schemas s ON t.schema_id=s.schema_id
    WHERE
        s.name=N'quik'
        AND t.name=N'money_limits'
) BEGIN
CREATE TABLE
    quik.money_limits (
        load_date date NOT NULL DEFAULT (getdate ()),
        client_code varchar(12) NOT NULL,
        ccy varchar(4) NOT NULL,
        position_code varchar(4) NOT NULL,
        settle_code varchar(5) NOT NULL,
        firm_code varchar(12) NOT NULL,
        firm_name varchar(128) NULL,
        balance DECIMAL(19,4) NULL,
        source_date date NOT NULL DEFAULT (getdate()),
        ts timestamp NOT NULL,

        CONSTRAINT PK_quik_money_limits PRIMARY KEY CLUSTERED (
            load_date ASC,
            client_code ASC,
            ccy ASC,
            position_code ASC,
            settle_code ASC,
            firm_code ASC
        ),
    );

END 
GO