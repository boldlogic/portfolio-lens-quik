package firm

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"go.uber.org/zap"
)

const (
	insertFirms = `
		INSERT INTO quik.firms (code, name)
		OUTPUT inserted.firm_id, inserted.code, inserted.name
		SELECT @p1, @p2
		WHERE NOT EXISTS (
			SELECT 1 FROM quik.firms WITH (UPDLOCK, HOLDLOCK) WHERE code = @p1
		)
	`

	updateFirmName = `
		UPDATE quik.firms
		SET name = @p2
		OUTPUT inserted.firm_id, inserted.code, inserted.name
		WHERE firm_id = @p1
`
)

func (r *FirmsRepo) InsertFirm(ctx context.Context, code string, name string) (quik.Firm, error) {
	var res quik.Firm
	r.repo.Logger.Debug("сохранение фирмы брокера", zap.String("code", code), zap.String("name", name))

	err := r.repo.Db.QueryRowContext(ctx, insertFirms, code, name).Scan(&res.Id, &res.Code, &res.Name)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.Firm{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.repo.Logger.Warn("фирма с таким кодом уже существует", zap.String("code", code))
			return quik.Firm{}, models.ErrConflict
		}
		r.repo.Logger.Error("ошибка сохранения фирмы брокера", zap.String("code", code), zap.String("name", name), zap.Error(err))
		return quik.Firm{}, models.ErrSavingData
	}

	r.repo.Logger.Debug("фирма брокера успешно сохранена", zap.Uint8("firm_id", res.Id), zap.String("code", res.Code), zap.String("name", res.Name))
	return res, nil
}

func (r *FirmsRepo) UpdateFirm(ctx context.Context, id uint8, name string) (quik.Firm, error) {
	var res quik.Firm
	row := r.repo.Db.QueryRowContext(ctx, updateFirmName, id, name)
	err := row.Scan(&res.Id, &res.Code, &res.Name)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.Firm{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.repo.Logger.Debug("фирма не найдена для обновления", zap.Uint8("id", id))
			return quik.Firm{}, models.ErrNotFound
		}
		r.repo.Logger.Error("ошибка обновления фирмы", zap.Uint8("firm_id", id), zap.Error(err))
		return quik.Firm{}, models.ErrSavingData
	}
	r.repo.Logger.Debug("фирма обновлена", zap.Uint8("firm_id", res.Id))
	return res, nil
}
