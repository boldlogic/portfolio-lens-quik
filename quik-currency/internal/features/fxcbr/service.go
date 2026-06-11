package fxcbr

import (
	"context"

	"go.uber.org/zap"
)

type service struct {
	logger *zap.Logger
	repo   fxRepo
}

func NewService(logger *zap.Logger, repo fxRepo) *service {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &service{
		logger: logger,
		repo:   repo,
	}
}

type fxRepo interface {
	MergeFxCBRRates(ctx context.Context) error
}

func (s *service) MergeFxCBRRatesQuik(ctx context.Context) error {
	if err := s.repo.MergeFxCBRRates(ctx); err != nil {
		s.logger.Error("произошла ошибка при сохранении кросс-курсов валют из QUIK", zap.Error(err))
		return err
	}
	return nil
}
