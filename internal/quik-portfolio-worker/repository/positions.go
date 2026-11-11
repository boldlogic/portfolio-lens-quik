package repository

import (
	"context"
	"time"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"go.uber.org/zap"
)

const listPositionLoadDatesSQL = `
	SELECT load_date
	FROM (
		SELECT load_date FROM quik.money_limits
		UNION
		SELECT load_date FROM quik.security_limits
		UNION
		SELECT load_date FROM quik.security_limits_otc
	) d
	ORDER BY load_date;
`

const deletePositionsWithoutLimitsSQL = `
	SET NOCOUNT ON;
	SET XACT_ABORT ON;
	BEGIN TRAN;
	DELETE p
	FROM portfolio.positions p
	WHERE NOT EXISTS (
		SELECT 1 FROM quik.money_limits m WHERE m.load_date = p.load_date
		UNION ALL
		SELECT 1 FROM quik.security_limits s WHERE s.load_date = p.load_date
		UNION ALL
		SELECT 1 FROM quik.security_limits_otc o WHERE o.load_date = p.load_date
	);
	COMMIT;
`

const replacePositionsForDateSQL = `
	SET NOCOUNT ON;
	SET XACT_ABORT ON;

	DECLARE @d date = CAST(@p1 AS date);

	DECLARE @money TABLE (
		client_id int NULL,
		instrument_id bigint NULL,
		balance decimal(19, 4) NULL
	);

	INSERT INTO @money (client_id, instrument_id, balance)
	SELECT
		c.client_id,
		i.instrument_id,
		CAST(ISNULL(p.balance, 0) AS decimal(19, 4))
	FROM (
		SELECT
			li.client_code,
			li.firm_code,
			li.currency_code,
			li.balance,
			li.settle_code,
			settle_max = MAX(li.settle_code) OVER (
				PARTITION BY li.load_date, li.client_code, li.firm_code, li.position_code, li.currency_code
			)
		FROM quik.money_limits li
		WHERE li.load_date = @d
	) p
	CROSS APPLY (
		SELECT ccy = CASE
			WHEN UPPER(LTRIM(RTRIM(p.currency_code))) IN ('SUR', 'RUR') THEN 'RUB'
			ELSE UPPER(LTRIM(RTRIM(p.currency_code)))
		END
	) norm
	LEFT JOIN ref.firms f ON f.code = p.firm_code
	LEFT JOIN ref.clients c ON c.firm_id = f.firm_id AND c.code = p.client_code
	LEFT JOIN ref.external_systems es ON es.ext_system = 'QUIK'
	LEFT JOIN ref.external_codes ec
		ON ec.ext_system_id = es.ext_system_id
		AND ec.ext_code_type_id = 1
		AND ec.ext_code = norm.ccy
	LEFT JOIN ref.currencies cur_ec
		ON ec.external_code_id IS NOT NULL
		AND cur_ec.iso_code = ec.internal_id
	LEFT JOIN ref.currencies cur_iso
		ON ec.external_code_id IS NULL
		AND cur_iso.iso_char_code = norm.ccy
	LEFT JOIN ref.instrument_types it ON it.title = N'Валюта'
	LEFT JOIN ref.instruments i
		ON i.type_id = it.type_id
		AND i.ticker = RTRIM(COALESCE(cur_ec.iso_char_code, cur_iso.iso_char_code))
	WHERE p.settle_code = p.settle_max;

	DECLARE @security TABLE (
		trade_account_id int NULL,
		instrument_id bigint NULL,
		balance decimal(19, 4) NULL
	);

	INSERT INTO @security (trade_account_id, instrument_id, balance)
	SELECT t.trade_account_id, i.instrument_id, p.balance
	FROM (
		SELECT
			li.client_code,
			li.firm_code,
			li.trade_account,
			li.sec_code,
			li.balance,
			li.settle_code,
			settle_max = MAX(li.settle_code) OVER (
				PARTITION BY li.load_date, li.client_code, li.firm_code, li.trade_account, li.sec_code
			)
		FROM quik.security_limits li
		WHERE li.load_date = @d
	) p
	LEFT JOIN ref.firms f ON f.code = p.firm_code
	LEFT JOIN ref.clients c ON c.firm_id = f.firm_id AND c.code = p.client_code
	LEFT JOIN ref.trade_accounts t
		ON t.client_id = c.client_id
		AND t.venue = 'exchange'
		AND t.code = p.trade_account
	LEFT JOIN ref.external_systems es ON es.ext_system = 'QUIK'
	LEFT JOIN ref.external_codes ec
		ON ec.ext_system_id = es.ext_system_id
		AND ec.ext_code_type_id = 3
		AND ec.ext_code = p.sec_code
	LEFT JOIN ref.instruments i ON i.instrument_id = ec.internal_id
	WHERE p.settle_code = p.settle_max;

	DECLARE @otc TABLE (
		trade_account_id int NULL,
		instrument_id bigint NULL,
		balance decimal(19, 4) NULL
	);

	INSERT INTO @otc (trade_account_id, instrument_id, balance)
	SELECT t.trade_account_id, i.instrument_id, p.balance
	FROM (
		SELECT
			li.client_code,
			li.firm_code,
			li.trade_account,
			li.sec_code,
			li.balance,
			li.settle_code,
			settle_max = MAX(li.settle_code) OVER (
				PARTITION BY li.load_date, li.client_code, li.firm_code, li.trade_account, li.sec_code
			)
		FROM quik.security_limits_otc li
		WHERE li.load_date = @d
	) p
	LEFT JOIN ref.firms f ON f.code = p.firm_code
	LEFT JOIN ref.clients c ON c.firm_id = f.firm_id AND c.code = p.client_code
	LEFT JOIN ref.trade_accounts t
		ON t.client_id = c.client_id
		AND t.venue = 'otc'
		AND t.code = p.trade_account
	LEFT JOIN ref.external_systems es ON es.ext_system = 'QUIK'
	LEFT JOIN ref.external_codes ec
		ON ec.ext_system_id = es.ext_system_id
		AND ec.ext_code_type_id = 3
		AND ec.ext_code = p.sec_code
	LEFT JOIN ref.instruments i ON i.instrument_id = ec.internal_id
	WHERE p.settle_code = p.settle_max;

	DECLARE @pos TABLE (
		client_id int NULL,
		trade_account_id int NULL,
		instrument_id bigint NOT NULL,
		balance decimal(19, 4) NULL
	);

	INSERT INTO @pos (client_id, trade_account_id, instrument_id, balance)
	SELECT client_id, NULL, instrument_id, SUM(balance)
	FROM @money
	WHERE client_id IS NOT NULL AND instrument_id IS NOT NULL
	GROUP BY client_id, instrument_id;

	INSERT INTO @pos (client_id, trade_account_id, instrument_id, balance)
	SELECT NULL, trade_account_id, instrument_id, SUM(balance)
	FROM @security
	WHERE trade_account_id IS NOT NULL AND instrument_id IS NOT NULL
	GROUP BY trade_account_id, instrument_id;

	INSERT INTO @pos (client_id, trade_account_id, instrument_id, balance)
	SELECT NULL, trade_account_id, instrument_id, SUM(balance)
	FROM @otc
	WHERE trade_account_id IS NOT NULL AND instrument_id IS NOT NULL
	GROUP BY trade_account_id, instrument_id;

	BEGIN TRAN;

	DELETE p
	FROM portfolio.positions AS p
	WHERE p.load_date = @d
		AND NOT EXISTS (
			SELECT 1
			FROM @pos AS s
			WHERE s.instrument_id = p.instrument_id
				AND EXISTS (SELECT p.client_id INTERSECT SELECT s.client_id)
				AND EXISTS (SELECT p.trade_account_id INTERSECT SELECT s.trade_account_id)
		);

	UPDATE p
	SET balance = s.balance
	FROM portfolio.positions AS p
	JOIN @pos AS s
		ON s.instrument_id = p.instrument_id
		AND EXISTS (SELECT p.client_id INTERSECT SELECT s.client_id)
		AND EXISTS (SELECT p.trade_account_id INTERSECT SELECT s.trade_account_id)
	WHERE p.load_date = @d
		AND EXISTS (SELECT p.balance EXCEPT SELECT s.balance);

	INSERT INTO portfolio.positions (load_date, client_id, trade_account_id, instrument_id, balance)
	SELECT @d, s.client_id, s.trade_account_id, s.instrument_id, s.balance
	FROM @pos AS s
	WHERE NOT EXISTS (
		SELECT 1
		FROM portfolio.positions AS p
		WHERE p.load_date = @d
			AND p.instrument_id = s.instrument_id
			AND EXISTS (SELECT p.client_id INTERSECT SELECT s.client_id)
			AND EXISTS (SELECT p.trade_account_id INTERSECT SELECT s.trade_account_id)
	);

	COMMIT;

	SELECT
		money_no_client = (SELECT COUNT(*) FROM @money WHERE client_id IS NULL),
		money_no_instrument = (SELECT COUNT(*) FROM @money WHERE client_id IS NOT NULL AND instrument_id IS NULL),
		security_no_account = (SELECT COUNT(*) FROM @security WHERE trade_account_id IS NULL),
		security_no_instrument = (SELECT COUNT(*) FROM @security WHERE trade_account_id IS NOT NULL AND instrument_id IS NULL),
		otc_no_account = (SELECT COUNT(*) FROM @otc WHERE trade_account_id IS NULL),
		otc_no_instrument = (SELECT COUNT(*) FROM @otc WHERE trade_account_id IS NOT NULL AND instrument_id IS NULL);
`

type PositionSkip struct {
	MoneyNoClient        int
	MoneyNoInstrument    int
	SecurityNoAccount    int
	SecurityNoInstrument int
	OtcNoAccount         int
	OtcNoInstrument      int
}

func (s PositionSkip) empty() bool {
	return s.MoneyNoClient == 0 &&
		s.MoneyNoInstrument == 0 &&
		s.SecurityNoAccount == 0 &&
		s.SecurityNoInstrument == 0 &&
		s.OtcNoAccount == 0 &&
		s.OtcNoInstrument == 0
}

func (r *Repository) ListPositionLoadDates(ctx context.Context) ([]time.Time, error) {
	rows, err := r.Db.QueryContext(ctx, listPositionLoadDatesSQL)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка чтения дат лимитов", zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}
			r.Logger.Error("ошибка чтения даты лимита", zap.Error(err))
			return nil, models.ErrRetrievingData
		}
		dates = append(dates, date)
	}
	if err := rows.Err(); err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка чтения дат лимитов", zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	return dates, nil
}

func (r *Repository) ReplacePositionsForDate(ctx context.Context, loadDate time.Time) error {
	var skip PositionSkip
	err := r.Db.QueryRowContext(ctx, replacePositionsForDateSQL, loadDate).Scan(
		&skip.MoneyNoClient,
		&skip.MoneyNoInstrument,
		&skip.SecurityNoAccount,
		&skip.SecurityNoInstrument,
		&skip.OtcNoAccount,
		&skip.OtcNoInstrument,
	)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.Logger.Error("ошибка сборки позиций", zap.Time("load_date", loadDate), zap.Error(err))
		return models.ErrSavingData
	}
	if !skip.empty() {
		r.Logger.Info("пропуск групп при сборке позиций",
			zap.Time("load_date", loadDate),
			zap.Int("money_no_client", skip.MoneyNoClient),
			zap.Int("money_no_instrument", skip.MoneyNoInstrument),
			zap.Int("security_no_account", skip.SecurityNoAccount),
			zap.Int("security_no_instrument", skip.SecurityNoInstrument),
			zap.Int("otc_no_account", skip.OtcNoAccount),
			zap.Int("otc_no_instrument", skip.OtcNoInstrument),
		)
	}
	return nil
}

func (r *Repository) DeletePositionsWithoutLimits(ctx context.Context) error {
	if _, err := r.Db.ExecContext(ctx, deletePositionsWithoutLimitsSQL); err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.Logger.Error("ошибка удаления позиций без лимитов", zap.Error(err))
		return models.ErrSavingData
	}
	return nil
}
