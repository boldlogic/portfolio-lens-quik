-- +goose Up

CREATE TABLE tinvest.subscriptions (
    subscription_id BIGINT        NOT NULL IDENTITY(1, 1),
    instrument_uid  NVARCHAR(64)  NOT NULL CONSTRAINT DF_tinvest_subscriptions_uid DEFAULT N'',
    ticker          NVARCHAR(64)  NOT NULL CONSTRAINT DF_tinvest_subscriptions_ticker DEFAULT N'',
    class_code      NVARCHAR(32)  NOT NULL CONSTRAINT DF_tinvest_subscriptions_class DEFAULT N'',
    figi            NVARCHAR(32)  NOT NULL CONSTRAINT DF_tinvest_subscriptions_figi DEFAULT N'',
    asset_uid       NVARCHAR(64)  NOT NULL CONSTRAINT DF_tinvest_subscriptions_asset DEFAULT N'',
    CONSTRAINT PK_tinvest_subscriptions PRIMARY KEY CLUSTERED (subscription_id)
);

CREATE UNIQUE INDEX UQ_tinvest_subscriptions_uid
    ON tinvest.subscriptions(instrument_uid)
    WHERE instrument_uid <> N'';

CREATE UNIQUE INDEX UQ_tinvest_subscriptions_ticker
    ON tinvest.subscriptions(ticker, class_code)
    WHERE instrument_uid = N'' AND ticker <> N'';

GRANT SELECT, INSERT, UPDATE ON OBJECT::tinvest.subscriptions TO t_invest_portfolio_writer;

-- +goose Down

DROP TABLE tinvest.subscriptions;
