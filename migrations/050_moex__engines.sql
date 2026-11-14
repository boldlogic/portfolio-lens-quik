-- +goose Up
-- +goose StatementBegin
CREATE TABLE moex.engines (
    engine_id INT IDENTITY(1,1) PRIMARY KEY,
    engine NVARCHAR(50) NOT NULL UNIQUE
);
CREATE INDEX IX_engines_engine ON moex.engines(engine);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS moex.engines;
