-- +goose Up
-- +goose StatementBegin
CREATE TABLE moex.securities (
    sec_id NVARCHAR(20) PRIMARY KEY,
    isin NVARCHAR(12),
    full_name NVARCHAR(255) NOT NULL,
    short_name NVARCHAR(50),
    reg_number NVARCHAR(50),
    type_code NVARCHAR(20),
    group_code NVARCHAR(50),
    issue_date DATE,
    face_unit NVARCHAR(10),
    face_value DECIMAL(18,6),
    initial_face_value DECIMAL(18,6),
    coupon_value DECIMAL(18,6),
    coupon_frequency INT,
    maturity_date DATE,
    CONSTRAINT FK_securities_security_types FOREIGN KEY (type_code) REFERENCES moex.security_types(type_code),
    CONSTRAINT FK_securities_security_groups FOREIGN KEY (group_code) REFERENCES moex.security_groups(group_code)
);
CREATE INDEX IX_securities_isin ON moex.securities(isin);
CREATE INDEX IX_securities_type_code ON moex.securities(type_code);
CREATE INDEX IX_securities_group_code ON moex.securities(group_code);
CREATE INDEX IX_securities_reg_number ON moex.securities(reg_number);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS moex.securities;
