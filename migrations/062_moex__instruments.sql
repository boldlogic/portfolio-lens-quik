-- +goose Up
-- +goose StatementBegin
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

-- +goose Down
DROP TABLE IF EXISTS moex.instruments;
