IF EXISTS (
    SELECT 1
    FROM sys.table_types tt
    INNER JOIN sys.schemas s ON tt.schema_id = s.schema_id
    WHERE s.name = N'api'
      AND tt.name = N'client_code_list'
)
BEGIN
    DROP TYPE api.client_code_list;
END;
GO

CREATE TYPE api.client_code_list AS TABLE (
    client_code varchar(12) NOT NULL
);
GO
