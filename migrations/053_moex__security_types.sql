-- +goose Up
-- +goose StatementBegin
CREATE TABLE moex.security_types (
    type_code NVARCHAR(20) PRIMARY KEY,
    type_name NVARCHAR(100) NOT NULL
);
CREATE INDEX IX_security_types_type_name ON moex.security_types(type_name);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS moex.security_types;
