package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/packages/shutdown"
	"github.com/boldlogic/quik-portfolio/pkg/models"
	"github.com/boldlogic/quik-portfolio/pkg/models/quik"
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
	selectFirms = `
		SELECT firm_id, code, name
		FROM quik.firms
		ORDER BY firm_id
`
	selectFirmByName = `
		SELECT  firm_id
			,code
			,name
		FROM quik.firms
		where name=@p1
`
	selectFirmByID = `
		SELECT firm_id, code, name
		FROM quik.firms
		WHERE firm_id = @p1
`
	updateFirmName = `
		UPDATE quik.firms
		SET name = @p2
		OUTPUT inserted.firm_id, inserted.code, inserted.name
		WHERE firm_id = @p1
`
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

func (r *Repository) SelectFirms(ctx context.Context) ([]quik.Firm, error) {
	rows, err := r.Db.QueryContext(ctx, selectFirms)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}
		r.Logger.Error("ошибка получения списка фирм", zap.Error(err))
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
			r.Logger.Error("ошибка чтения фирмы", zap.Error(err))
			return nil, models.ErrRetrievingData
		}
		result = append(result, f)
	}
	if rows.Err() != nil {
		r.Logger.Error("ошибка получения списка фирм", zap.Error(rows.Err()))

		return nil, models.ErrRetrievingData
	}
	return result, nil
}

func (r *Repository) SelectFirmByID(ctx context.Context, id uint8) (quik.Firm, error) {
	var res quik.Firm
	row := r.Db.QueryRowContext(ctx, selectFirmByID, id)
	err := row.Scan(&res.Id, &res.Code, &res.Name)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.Firm{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Warn("фирма не найдена", zap.Uint8("firm_id", id))
			return quik.Firm{}, models.ErrNotFound
		}
		r.Logger.Error("ошибка получения фирмы", zap.Uint8("firm_id", id), zap.Error(err))
		return quik.Firm{}, models.ErrRetrievingData
	}
	return res, nil
}

func (r *Repository) UpdateFirm(ctx context.Context, id uint8, name string) (quik.Firm, error) {
	var res quik.Firm
	row := r.Db.QueryRowContext(ctx, updateFirmName, id, name)
	err := row.Scan(&res.Id, &res.Code, &res.Name)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.Firm{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Debug("фирма не найдена для обновления", zap.Uint8("id", id))
			return quik.Firm{}, models.ErrNotFound
		}
		r.Logger.Error("ошибка обновления фирмы", zap.Uint8("firm_id", id), zap.Error(err))
		return quik.Firm{}, models.ErrSavingData
	}
	r.Logger.Debug("фирма обновлена", zap.Uint8("firm_id", res.Id))
	return res, nil
}

func (r *Repository) SelectFirmByName(ctx context.Context, name string) (quik.Firm, error) {
	var res quik.Firm
	r.Logger.Debug("получение фирмы по имени", zap.String("name", name))
	row := r.Db.QueryRowContext(ctx, selectFirmByName, name)
	err := row.Scan(&res.Id, &res.Code, &res.Name)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.Firm{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Warn("фирма не найдена", zap.String("name", name))
			return quik.Firm{}, models.ErrNotFound
		}
		r.Logger.Error("ошибка получения фирмы по имени", zap.String("name", name), zap.Error(err))
		return quik.Firm{}, models.ErrRetrievingData
	}
	return res, nil
}

func (r *Repository) InsertFirm(ctx context.Context, code string, name string) (quik.Firm, error) {
	var res quik.Firm
	r.Logger.Debug("сохранение фирмы брокера", zap.String("code", code), zap.String("name", name))

	err := r.Db.QueryRowContext(ctx, insertFirms, code, name).Scan(&res.Id, &res.Code, &res.Name)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return quik.Firm{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Warn("фирма с таким кодом уже существует", zap.String("code", code))
			return quik.Firm{}, models.ErrConflict
		}
		r.Logger.Error("ошибка сохранения фирмы брокера", zap.String("code", code), zap.String("name", name), zap.Error(err))
		return quik.Firm{}, models.ErrSavingData
	}

	r.Logger.Debug("фирма брокера успешно сохранена", zap.Uint8("firm_id", res.Id), zap.String("code", res.Code), zap.String("name", res.Name))
	return res, nil
}

func (r *Repository) SyncFirmsFromLimits(ctx context.Context) error {
	_, err := r.Db.ExecContext(ctx, syncFirmsFromLimits)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.Logger.Error("ошибка синхронизации фирм из лимитов", zap.Error(err))
		return models.ErrSavingData
	}
	r.Logger.Debug("фирмы из лимитов синхронизированы")
	return nil
}
