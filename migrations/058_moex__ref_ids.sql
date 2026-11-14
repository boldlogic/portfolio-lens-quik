-- +goose Up
-- +goose StatementBegin
ALTER TABLE moex.security_types ADD iss_type_id INT NULL;
ALTER TABLE moex.security_groups ADD iss_group_id INT NULL;
ALTER TABLE moex.engines ADD iss_engine_id INT NULL;
CREATE INDEX IX_security_types_iss_type_id ON moex.security_types(iss_type_id);
CREATE INDEX IX_security_groups_iss_group_id ON moex.security_groups(iss_group_id);
CREATE INDEX IX_engines_iss_engine_id ON moex.engines(iss_engine_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS IX_engines_iss_engine_id ON moex.engines;
DROP INDEX IF EXISTS IX_security_groups_iss_group_id ON moex.security_groups;
DROP INDEX IF EXISTS IX_security_types_iss_type_id ON moex.security_types;
ALTER TABLE moex.engines DROP COLUMN iss_engine_id;
ALTER TABLE moex.security_groups DROP COLUMN iss_group_id;
ALTER TABLE moex.security_types DROP COLUMN iss_type_id;
-- +goose StatementEnd
