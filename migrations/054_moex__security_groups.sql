-- +goose Up
-- +goose StatementBegin
CREATE TABLE moex.security_groups (
    group_code NVARCHAR(50) PRIMARY KEY,
    group_name NVARCHAR(100) NOT NULL
);
CREATE INDEX IX_security_groups_group_name ON moex.security_groups(group_name);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS moex.security_groups;
