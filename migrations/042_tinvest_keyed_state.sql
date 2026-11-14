-- +goose Up
-- +goose StatementBegin
IF SCHEMA_ID(N'tinvest') IS NULL EXEC(N'CREATE SCHEMA tinvest');

IF OBJECT_ID(N'tinvest.current_state', N'U') IS NULL
BEGIN
    CREATE TABLE tinvest.current_state (
        entity NVARCHAR(64) NOT NULL,
        scope_key NVARCHAR(128) NOT NULL,
        entity_key NVARCHAR(128) NOT NULL,
        observed_at DATETIME2(7) NOT NULL,
        payload NVARCHAR(MAX) NOT NULL,
        CONSTRAINT PK_tinvest_current_state PRIMARY KEY (entity, scope_key, entity_key)
    );
END;

IF OBJECT_ID(N'tinvest.history_record', N'U') IS NULL
BEGIN
    CREATE TABLE tinvest.history_record (
        entity NVARCHAR(64) NOT NULL,
        scope_key NVARCHAR(128) NOT NULL,
        record_key NVARCHAR(128) NOT NULL,
        event_at DATETIME2(7) NULL,
        payload NVARCHAR(MAX) NOT NULL,
        CONSTRAINT PK_tinvest_history_record PRIMARY KEY (entity, scope_key, record_key)
    );
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
IF OBJECT_ID(N'tinvest.history_record', N'U') IS NOT NULL DROP TABLE tinvest.history_record;
IF OBJECT_ID(N'tinvest.current_state', N'U') IS NOT NULL DROP TABLE tinvest.current_state;
-- +goose StatementEnd
