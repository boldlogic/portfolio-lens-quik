-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS moex.instruments;

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
    s.face_unit,
    s.initial_face_value,
    s.face_value,
    s.coupon_percent,
    s.coupon_value,
    s.coupon_frequency,
    s.maturity_date
FROM moex.securities s
LEFT JOIN moex.security_types st ON st.type_code = s.type_code
LEFT JOIN moex.security_groups sg ON sg.group_code = s.group_code;
');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS moex.instruments;

CREATE TABLE moex.instruments (
    sec_id NVARCHAR(20) PRIMARY KEY,
    full_name NVARCHAR(255) NOT NULL,
    isin NVARCHAR(12) NULL,
    short_name NVARCHAR(50) NULL,
    reg_number NVARCHAR(50) NULL,
    lat_name NVARCHAR(250) NULL,
    currency_id NVARCHAR(10) NULL,
    type_code NVARCHAR(20) NULL,
    type_name NVARCHAR(100) NULL,
    group_code NVARCHAR(50) NULL,
    group_name NVARCHAR(100) NULL,
    iss_type_id INT NULL,
    iss_group_id INT NULL,
    issue_date DATE NULL,
    lot_size INT NULL,
    face_unit NVARCHAR(10) NULL,
    initial_face_value DECIMAL(18,6) NULL,
    face_value DECIMAL(18,6) NULL,
    coupon_percent DECIMAL(18,6) NULL,
    coupon_value DECIMAL(18,6) NULL,
    coupon_frequency INT NULL,
    maturity_date DATE NULL
);
CREATE INDEX IX_instruments_type_code ON moex.instruments(type_code);
CREATE INDEX IX_instruments_group_code ON moex.instruments(group_code);
-- +goose StatementEnd
