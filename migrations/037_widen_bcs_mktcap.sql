-- +goose Up

ALTER TABLE bcs.security_estimates ALTER COLUMN mktcap DECIMAL(28, 8) NULL;

-- +goose Down

ALTER TABLE bcs.security_estimates ALTER COLUMN mktcap DECIMAL(19, 8) NULL;
