package client

import (
	"context"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"go.uber.org/zap"
)

const (
	syncClientsFromLimits = `
		INSERT INTO ref.clients (firm_id, code)
		SELECT DISTINCT f.firm_id, src.client_code
		FROM (
			SELECT firm_code, client_code
			FROM quik.money_limits
			WHERE LTRIM(RTRIM(client_code)) <> ''
				AND LTRIM(RTRIM(firm_code)) <> ''
			UNION
			SELECT firm_code, client_code
			FROM quik.security_limits
			WHERE LTRIM(RTRIM(client_code)) <> ''
				AND LTRIM(RTRIM(firm_code)) <> ''
			UNION
			SELECT firm_code, client_code
			FROM quik.security_limits_otc
			WHERE LTRIM(RTRIM(client_code)) <> ''
				AND LTRIM(RTRIM(firm_code)) <> ''
		) src
		JOIN ref.firms f ON f.code = src.firm_code
		WHERE NOT EXISTS (
			SELECT 1
			FROM ref.clients c
			WHERE c.firm_id = f.firm_id
				AND c.code = src.client_code
		);
	`

	countSkippedClients = `
		SELECT
			(SELECT COUNT(*) FROM quik.money_limits
				WHERE LTRIM(RTRIM(ISNULL(client_code, ''))) = ''
					OR LTRIM(RTRIM(ISNULL(firm_code, ''))) = '')
			+ (SELECT COUNT(*) FROM quik.security_limits
				WHERE LTRIM(RTRIM(ISNULL(client_code, ''))) = ''
					OR LTRIM(RTRIM(ISNULL(firm_code, ''))) = '')
			+ (SELECT COUNT(*) FROM quik.security_limits_otc
				WHERE LTRIM(RTRIM(ISNULL(client_code, ''))) = ''
					OR LTRIM(RTRIM(ISNULL(firm_code, ''))) = '');
	`
)

func (r *ClientsRepo) SyncClientsFromLimits(ctx context.Context) error {
	if _, err := r.repo.Db.ExecContext(ctx, syncClientsFromLimits); err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.repo.Logger.Error("ошибка синхронизации клиентов из лимитов", zap.Error(err))
		return models.ErrSavingData
	}

	var skipped int
	if err := r.repo.Db.QueryRowContext(ctx, countSkippedClients).Scan(&skipped); err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.repo.Logger.Error("ошибка подсчета лимитов без кода клиента или фирмы", zap.Error(err))
		return models.ErrRetrievingData
	}
	if skipped > 0 {
		r.repo.Logger.Info("пропущены лимиты без кода клиента или фирмы", zap.Int("skipped", skipped))
	}
	return nil
}
