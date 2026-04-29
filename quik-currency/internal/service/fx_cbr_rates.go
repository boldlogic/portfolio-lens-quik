package service

import (
	"context"

	"go.uber.org/zap"
)

func (s *Service) MergeFxCBRRatesQuik(ctx context.Context) error {
	if err := s.currencyRepo.MergeFxCBRRatesQuik(ctx); err != nil {
		s.logger.Error("произошла ошибка при сохранении кросс-курсов валют из QUIK", zap.Error(err))
		return err
	}

	return nil
}
