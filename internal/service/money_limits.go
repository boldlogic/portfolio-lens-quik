package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/boldlogic/quik-portfolio/internal/models"
	"github.com/boldlogic/quik-portfolio/pkg/currencies"
	md "github.com/boldlogic/quik-portfolio/pkg/models"
)

func (s *Service) GetMoneyLimits(ctx context.Context, date time.Time) ([]models.MoneyLimit, error) {
	maxDate, err := s.repo.SelectMoneyLimitsMaxDate(ctx)
	if err != nil {
		return nil, err
	}
	if maxDate == nil {
		return []models.MoneyLimit{}, nil
	}
	if err := checkLimitDate(date, *maxDate); err != nil {
		return nil, err
	}
	return s.repo.SelectMoneyLimits(ctx, date)
}

func (s *Service) CreateMoneyLimit(ctx context.Context, ml models.MoneyLimit) (models.MoneyLimit, error) {

	maxDate, err := s.repo.SelectMoneyLimitsMaxDate(ctx)
	if err != nil && !errors.Is(err, md.ErrNotFound) {
		return models.MoneyLimit{}, err
	}
	if err := checkLimitDate(ml.LoadDate, minRollForwardDate(maxDate)); err != nil {
		return models.MoneyLimit{}, err
	}

	normCcy := normalizeQuikCcy(ml.Currency)
	if err := currencies.CheckCurrencyCode(normCcy); err != nil {
		return models.MoneyLimit{}, fmt.Errorf("%w: %s", md.ErrBusinessValidation, err.Error())
	}

	if ml.SettleCode == "" {
		ml.SettleCode = md.SettleCodeTx
	}

	err = ml.SettleCode.Validate()
	if err != nil {
		return models.MoneyLimit{}, fmt.Errorf("%w: %s", md.ErrBusinessValidation, err.Error())
	}

	created, err := s.repo.InsertMoneyLimit(ctx, ml)
	if err != nil {
		if errors.Is(err, md.ErrNotFound) {
			return models.MoneyLimit{}, fmt.Errorf("%w: некорректное имя фирмы", md.ErrBusinessValidation)
		}
		if errors.Is(err, md.ErrConflict) {
			return models.MoneyLimit{}, fmt.Errorf("%w: load_date=%s client_code=%s ccy=%s position_code=%s settle_code=%s firm_name=%s",
				md.ErrConflict,
				ml.LoadDate.Format(md.ISODateFormat),
				ml.ClientCode,
				ml.Currency,
				ml.PositionCode,
				ml.SettleCode,
				ml.FirmName)
		}

		return models.MoneyLimit{}, err
	}
	return created, nil
}

func (s *Service) DoRollForwardMoneyLimits(ctx context.Context) error {
	return doRollForward(ctx,
		s.repo.SelectMoneyLimitsMaxDate,
		s.repo.InsertMoneyLimitsCopy,
		s.repo.DeleteMoneyLimits,
	)
}
