-- +goose Up

DROP TABLE bcs.limit_payload;
DROP TABLE bcs.portfolio;

-- +goose Down

CREATE TABLE bcs.portfolio (
    portfolio_id BIGINT            NOT NULL,
    received_at  DATETIMEOFFSET(7) NOT NULL,
    http_status  INT               NOT NULL,
    raw_body     VARBINARY(MAX)    NOT NULL,
    CONSTRAINT PK_bcs_portfolio PRIMARY KEY CLUSTERED (portfolio_id)
);

CREATE TABLE bcs.limit_payload (
    payload_id BIGINT         NOT NULL,
    raw_body   VARBINARY(MAX) NOT NULL,
    CONSTRAINT PK_bcs_limit_payload PRIMARY KEY CLUSTERED (payload_id)
);
