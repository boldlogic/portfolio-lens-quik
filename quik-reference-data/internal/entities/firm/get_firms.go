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
	selectFirms = `
	SELECT firm_id, code, name
	FROM dbo.firms
	ORDER BY firm_id
`
	selectFirmByName = `
	SELECT  firm_id
		,code
		,name
	FROM dbo.firms
	where name=@p1
`
	selectFirmByID = `
	SELECT firm_id, code, name
	FROM dbo.firms
	WHERE firm_id = @p1
`
)

func (r *FirmsRepo) SelectFirms(ctx context.Context) ([]quik.Firm, error) {
	rows, err := r.repo.Db.QueryContext(ctx, selectFirms)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.repo.Logger.Error("ошибка получения списка фирм", zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	defer rows.Close()

	var result []quik.Firm
	for rows.Next() {
		var f quik.Firm
		if err := rows.Scan(&f.Id, &f.Code, &f.Name); err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}
			r.repo.Logger.Error("ошибка чтения фирмы", zap.Error(err))
			return nil, models.ErrRetrievingData
		}
		result = append(result, f)
	}
	if rows.Err() != nil {
		r.repo.Logger.Error("ошибка получения списка фирм", zap.Error(rows.Err()))

		return nil, models.ErrRetrievingData
	}
	return result, nil
}

func (r *FirmsRepo) SelectFirmByID(ctx context.Context, id uint8) (quik.Firm, error) {
	var res quik.Firm
	row := r.repo.Db.QueryRowContext(ctx, selectFirmByID, id)
	err := row.Scan(&res.Id, &res.Code, &res.Name)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.Firm{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.repo.Logger.Warn("фирма не найдена", zap.Uint8("firm_id", id))
			return quik.Firm{}, models.ErrNotFound
		}
		r.repo.Logger.Error("ошибка получения фирмы", zap.Uint8("firm_id", id), zap.Error(err))
		return quik.Firm{}, models.ErrRetrievingData
	}
	return res, nil
}

func (r *FirmsRepo) SelectFirmByName(ctx context.Context, name string) (quik.Firm, error) {
	var res quik.Firm
	r.repo.Logger.Debug("получение фирмы по имени", zap.String("name", name))
	row := r.repo.Db.QueryRowContext(ctx, selectFirmByName, name)
	err := row.Scan(&res.Id, &res.Code, &res.Name)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.Firm{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.repo.Logger.Warn("фирма не найдена", zap.String("name", name))
			return quik.Firm{}, models.ErrNotFound
		}
		r.repo.Logger.Error("ошибка получения фирмы по имени", zap.String("name", name), zap.Error(err))
		return quik.Firm{}, models.ErrRetrievingData
	}
	return res, nil
}
