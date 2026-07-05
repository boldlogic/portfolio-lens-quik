package service

import (
	"context"

	"github.com/boldlogic/packages/shutdown"
	"go.uber.org/zap"
)

func (s *Service) ImportQuikCrossRates(ctx context.Context) error {
	if err := s.repo.MergeFxCBRRates(ctx); err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		s.logger.Error("произошла ошибка при сохранении кросс-курсов валют из QUIK", zap.Error(err))
		return err
	}
	return nil
}
