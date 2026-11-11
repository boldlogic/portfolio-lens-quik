package instrumentcode

import (
	"context"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"go.uber.org/zap"
)

const syncInstrumentCodes = `
	INSERT INTO ref.external_codes (ext_system_id, ext_code, ext_code_type_id, internal_id)
	SELECT es.ext_system_id, i.ticker, 3, i.instrument_id
	FROM ref.instruments i
	JOIN (
		SELECT ticker
		FROM ref.instruments
		GROUP BY ticker
		HAVING COUNT(*) = 1
	) one ON one.ticker = i.ticker
	JOIN ref.external_systems es ON es.ext_system = 'QUIK'
	WHERE EXISTS (
			SELECT 1
			FROM quik.security_limits sl
			WHERE sl.sec_code = i.ticker
			UNION
			SELECT 1
			FROM quik.security_limits_otc ot
			WHERE ot.sec_code = i.ticker
		)
		AND NOT EXISTS (
			SELECT 1
			FROM ref.external_codes ec
			WHERE ec.ext_system_id = es.ext_system_id
				AND ec.ext_code_type_id = 3
				AND ec.ext_code = i.ticker
		);
`

func (r *InstrumentCodesRepo) SyncInstrumentCodes(ctx context.Context) error {
	if _, err := r.repo.Db.ExecContext(ctx, syncInstrumentCodes); err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.repo.Logger.Error("ошибка синхронизации кодов инструментов", zap.Error(err))
		return models.ErrSavingData
	}
	return nil
}
