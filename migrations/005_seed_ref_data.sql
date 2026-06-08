-- +goose Up

SET IDENTITY_INSERT ref.external_systems ON;

INSERT INTO ref.external_systems (ext_system_id, ext_system)
VALUES (1, 'QUIK');

SET IDENTITY_INSERT ref.external_systems OFF;

INSERT INTO ref.external_codes (ext_system_id, ext_code, ext_code_type_id, internal_id)
VALUES
    (1, N'GLD', 1, 959),
    (1, N'SLV', 1, 961),
    (1, N'PLT', 1, 962),
    (1, N'PLD', 1, 964);
