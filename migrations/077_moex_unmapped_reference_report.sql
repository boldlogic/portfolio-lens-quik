-- +goose Up

CREATE TABLE moex.common_sync_unmapped (
    unmapped_id BIGINT IDENTITY(1,1) NOT NULL,
    entity_type NVARCHAR(30) NOT NULL,
    source_code NVARCHAR(100) NOT NULL,
    context_code NVARCHAR(100) NOT NULL,
    first_seen DATETIME2 NOT NULL CONSTRAINT DF_moex_common_unmapped_first_seen DEFAULT SYSUTCDATETIME(),
    last_seen DATETIME2 NOT NULL CONSTRAINT DF_moex_common_unmapped_last_seen DEFAULT SYSUTCDATETIME(),
    CONSTRAINT PK_moex_common_sync_unmapped PRIMARY KEY CLUSTERED (unmapped_id),
    CONSTRAINT UQ_moex_common_sync_unmapped UNIQUE (entity_type, source_code, context_code)
);

-- +goose Down

THROW 51000, N'Отчёт несопоставленных кодов содержит диагностические данные и не откатывается автоматически.', 1;
