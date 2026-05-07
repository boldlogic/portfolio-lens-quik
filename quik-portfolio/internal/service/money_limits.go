package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-currency/pkg/currencies"
	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (s *Service) GetMoneyLimits(ctx context.Context, date time.Time) ([]quik.MoneyLimit, error) {
	maxDate, err := s.repo.SelectMoneyLimitsMaxDate(ctx)
	if err != nil {
		return nil, err
	}
	if maxDate == nil {
		return []quik.MoneyLimit{}, nil
	}
	if err := checkLimitDate(date, *maxDate); err != nil {
		return nil, err
	}
	return s.repo.SelectMoneyLimits(ctx, date)
}

// CreateMoneyLimit
func (s *Service) CreateMoneyLimit(ctx context.Context, ml quik.MoneyLimit) (quik.MoneyLimit, error) {

	maxDate, err := s.repo.SelectMoneyLimitsMaxDate(ctx)
	if err != nil && !errors.Is(err, md.ErrNotFound) {
		return quik.MoneyLimit{}, err
	}
	if err := checkLimitDate(ml.LoadDate, minRollForwardDate(maxDate)); err != nil {
		return quik.MoneyLimit{}, err
	}

	normCcy := normalizeQuikCcy(ml.Currency)
	if err := currencies.CheckCurrencyCode(normCcy); err != nil {
		return quik.MoneyLimit{}, fmt.Errorf("%w: %s", md.ErrBusinessValidation, err.Error())
	}

	if ml.SettleCode == "" {
		ml.SettleCode = quik.SettleCodeTx
	}

	err = ml.SettleCode.Validate()
	if err != nil {
		return quik.MoneyLimit{}, fmt.Errorf("%w: %s", md.ErrBusinessValidation, err.Error())
	}

	created, err := s.repo.InsertMoneyLimit(ctx, ml)
	if err != nil {
		if errors.Is(err, md.ErrNotFound) {
			return quik.MoneyLimit{}, fmt.Errorf("%w: некорректный код фирмы %s", md.ErrBusinessValidation, ml.FirmCode)
		}
		if errors.Is(err, md.ErrConflict) {
			return quik.MoneyLimit{}, fmt.Errorf("%w: loadDate=%s clientCode=%s currency=%s positionCode=%s settleCode=%s firmCode=%s",
				md.ErrConflict,
				ml.LoadDate.Format(dates.ISODateFormat),
				ml.ClientCode,
				ml.Currency,
				ml.PositionCode,
				ml.SettleCode,
				ml.FirmCode)
		}

		return quik.MoneyLimit{}, err
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
