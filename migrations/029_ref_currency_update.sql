-- +goose Up

ALTER TABLE ref.currencies add rw rowversion;
