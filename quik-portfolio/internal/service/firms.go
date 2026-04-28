package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (s *Service) GetFirms(ctx context.Context) ([]quik.Firm, error) {
	return s.repo.SelectFirms(ctx)
}

func (s *Service) GetFirmByID(ctx context.Context, id uint8) (quik.Firm, error) {
	firm, err := s.repo.SelectFirmByID(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return quik.Firm{}, fmt.Errorf("%w: фирма с id %d не найдена", models.ErrNotFound, id)
		}
		return quik.Firm{}, err
	}
	return firm, nil
}

func (s *Service) CreateFirm(ctx context.Context, code string, name string) (quik.Firm, error) {
	firm, err := s.repo.InsertFirm(ctx, code, name)
	if err != nil {
		if errors.Is(err, models.ErrConflict) {
			return quik.Firm{}, fmt.Errorf("%w: фирма c firmCode %s уже существует", models.ErrConflict, code)
		}
		return quik.Firm{}, err
	}
	return firm, nil
}

func (s *Service) UpdateFirm(ctx context.Context, id uint8, name string) (quik.Firm, error) {
	firm, err := s.repo.UpdateFirm(ctx, id, name)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return quik.Firm{}, fmt.Errorf("%w: фирма с id %d не найдена", models.ErrNotFound, id)
		}
		return quik.Firm{}, err
	}
	return firm, nil
}

func (s *Service) SyncFirmsFromLimits(ctx context.Context) error {
	return s.repo.SyncFirmsFromLimits(ctx)
}
