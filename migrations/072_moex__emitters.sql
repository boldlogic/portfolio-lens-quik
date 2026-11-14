-- +goose Up
-- +goose StatementBegin
CREATE TABLE moex.emitters (
    emitent_id INT NOT NULL PRIMARY KEY,
    title NVARCHAR(765) NOT NULL,
    inn NVARCHAR(30) NULL,
    okpo NVARCHAR(24) NULL
);
CREATE INDEX IX_emitters_inn ON moex.emitters(inn);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO moex.emitters (emitent_id, title)
SELECT DISTINCT TRY_CAST(s.emitter_id AS INT), CAST(TRY_CAST(s.emitter_id AS INT) AS NVARCHAR(50))
FROM moex.securities s
WHERE TRY_CAST(s.emitter_id AS INT) IS NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE moex.securities ADD emitent_id INT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE moex.securities SET emitent_id = TRY_CAST(emitter_id AS INT);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE moex.securities DROP COLUMN emitter_id;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE moex.securities ADD CONSTRAINT FK_securities_emitters
    FOREIGN KEY (emitent_id) REFERENCES moex.emitters(emitent_id);
-- +goose StatementEnd

-- +goose StatementBegin
DROP VIEW IF EXISTS moex.instruments;
-- +goose StatementEnd

-- +goose StatementBegin
EXEC(N'
CREATE VIEW moex.instruments AS
SELECT
    s.sec_id,
    s.full_name,
    s.isin,
    s.short_name,
    s.reg_number,
    s.lat_name,
    s.currency_id,
    s.type_code,
    st.type_name,
    s.group_code,
    sg.group_name,
    st.iss_type_id,
    sg.iss_group_id,
    s.issue_date,
    s.lot_size,
    b.face_unit,
    b.initial_face_value,
    b.face_value,
    b.coupon_percent,
    b.coupon_value,
    b.coupon_frequency,
    b.maturity_date,
    s.emitent_id,
    e.title AS emitent_title,
    e.inn AS emitent_inn,
    e.okpo AS emitent_okpo
FROM moex.securities s
LEFT JOIN moex.security_types st ON st.type_code = s.type_code
LEFT JOIN moex.security_groups sg ON sg.group_code = s.group_code
LEFT JOIN moex.security_bonds b ON b.sec_id = s.sec_id
LEFT JOIN moex.emitters e ON e.emitent_id = s.emitent_id;
');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS moex.instruments;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE moex.securities DROP CONSTRAINT FK_securities_emitters;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE moex.securities ADD emitter_id NVARCHAR(50) NULL;
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE moex.securities SET emitter_id = CAST(emitent_id AS NVARCHAR(50)) WHERE emitent_id IS NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE moex.securities DROP COLUMN emitent_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS moex.emitters;
-- +goose StatementEnd

-- +goose StatementBegin
EXEC(N'
CREATE VIEW moex.instruments AS
SELECT
    s.sec_id,
    s.full_name,
    s.isin,
    s.short_name,
    s.reg_number,
    s.lat_name,
    s.currency_id,
    s.type_code,
    st.type_name,
    s.group_code,
    sg.group_name,
    st.iss_type_id,
    sg.iss_group_id,
    s.issue_date,
    s.lot_size,
    b.face_unit,
    b.initial_face_value,
    b.face_value,
    b.coupon_percent,
    b.coupon_value,
    b.coupon_frequency,
    b.maturity_date
FROM moex.securities s
LEFT JOIN moex.security_types st ON st.type_code = s.type_code
LEFT JOIN moex.security_groups sg ON sg.group_code = s.group_code
LEFT JOIN moex.security_bonds b ON b.sec_id = s.sec_id;
');
-- +goose StatementEnd
