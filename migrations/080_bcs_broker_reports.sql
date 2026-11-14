-- +goose Up

CREATE TABLE bcs.report_imports (
    import_id       BIGINT            NOT NULL IDENTITY(1, 1),
    content_sha256  CHAR(64)          NOT NULL,
    file_name       NVARCHAR(260)     NOT NULL,
    file_format     NVARCHAR(8)       NOT NULL,
    import_status   NVARCHAR(16)      NOT NULL,
    row_count       INT               NOT NULL CONSTRAINT DF_bcs_report_imports_row_count DEFAULT (0),
    trade_count     INT               NOT NULL CONSTRAINT DF_bcs_report_imports_trade_count DEFAULT (0),
    operation_count INT               NOT NULL CONSTRAINT DF_bcs_report_imports_operation_count DEFAULT (0),
    error_text      NVARCHAR(2000)    NULL,
    created_at      DATETIMEOFFSET(7) NOT NULL CONSTRAINT DF_bcs_report_imports_created_at DEFAULT (SYSDATETIMEOFFSET()),
    finished_at     DATETIMEOFFSET(7) NULL,
    CONSTRAINT PK_bcs_report_imports PRIMARY KEY CLUSTERED (import_id),
    CONSTRAINT UQ_bcs_report_imports_sha256 UNIQUE (content_sha256),
    CONSTRAINT CK_bcs_report_imports_format CHECK (file_format IN (N'xls', N'xlsx')),
    CONSTRAINT CK_bcs_report_imports_status CHECK (import_status IN (N'received', N'processing', N'failed', N'completed')),
    CONSTRAINT CK_bcs_report_imports_counts CHECK (row_count >= 0 AND trade_count >= 0 AND operation_count >= 0)
);

CREATE TABLE bcs.report_rows (
    report_row_id BIGINT        NOT NULL IDENTITY(1, 1),
    import_id     BIGINT        NOT NULL,
    sheet_name    NVARCHAR(128) NOT NULL,
    row_number    INT           NOT NULL,
    raw_values    NVARCHAR(MAX) NOT NULL,
    CONSTRAINT PK_bcs_report_rows PRIMARY KEY CLUSTERED (report_row_id),
    CONSTRAINT FK_bcs_report_rows_import FOREIGN KEY (import_id)
        REFERENCES bcs.report_imports (import_id),
    CONSTRAINT UQ_bcs_report_rows_import_position UNIQUE (import_id, sheet_name, row_number),
    CONSTRAINT CK_bcs_report_rows_number CHECK (row_number > 0),
    CONSTRAINT CK_bcs_report_rows_json CHECK (ISJSON(raw_values) = 1)
);

CREATE INDEX IX_bcs_report_rows_import ON bcs.report_rows (import_id, row_number);

-- +goose Down

DROP TABLE bcs.report_rows;
DROP TABLE bcs.report_imports;
