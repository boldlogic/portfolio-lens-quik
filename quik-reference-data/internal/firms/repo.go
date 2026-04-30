package firms

import (
	"context"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/quik-reference-data/internal/infra"
	"go.uber.org/zap"
)

type FirmsRepo struct {
	repo *infra.Repository
}

func NewFirmsRepo(r *infra.Repository) *FirmsRepo {

	return &FirmsRepo{
		repo: r,
	}
}

const (
	syncFirmsFromLimits = `
		INSERT INTO quik.firms (code, name)
		SELECT DISTINCT src.code, src.name
		FROM (
			SELECT firm_code AS code, LTRIM(RTRIM(firm_name)) AS name
			FROM quik.money_limits
			WHERE firm_code IS NOT NULL AND LTRIM(RTRIM(firm_name)) <> ''
			UNION
			SELECT firm_code AS code, LTRIM(RTRIM(firm_name)) AS name
			FROM quik.security_limits
			WHERE firm_code IS NOT NULL AND LTRIM(RTRIM(firm_name)) <> ''
		) src
		WHERE NOT EXISTS (SELECT 1 FROM quik.firms f WHERE f.code = src.code);
	`
)

func (r *FirmsRepo) SyncFirmsFromLimits(ctx context.Context) error {
	_, err := r.repo.Db.ExecContext(ctx, syncFirmsFromLimits)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.repo.Logger.Error("ошибка синхронизации фирм из лимитов", zap.Error(err))
		return models.ErrSavingData
	}

	return nil
}
