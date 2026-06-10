-- +goose Up
ALTER TABLE quik.current_quotes 
ADD crossrate decimal(19, 8) NULL;

-- +goose Down
ALTER TABLE quik.current_quotes 
DROP COLUMN crossrate;