-- +goose Up
-- +goose StatementBegin
CREATE TABLE moex.sync_runs (
    run_id      BIGINT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    worker      NVARCHAR(50) NOT NULL,
    job         NVARCHAR(50) NOT NULL,
    started_at  DATETIME2 NOT NULL CONSTRAINT DF_sync_runs_started DEFAULT SYSUTCDATETIME(),
    finished_at DATETIME2 NULL,
    status      NVARCHAR(20) NOT NULL,
    error_msg   NVARCHAR(500) NULL
);
CREATE INDEX IX_sync_runs_worker_started ON moex.sync_runs(worker, started_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE moex.sync_runs;
-- +goose StatementEnd
