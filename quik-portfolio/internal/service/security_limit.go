package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (s *Service) GetSecurityLimits(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
	maxDate, err := s.repo.SelectSecurityLimitsMaxDate(ctx)
	if err != nil {
		return nil, err
	}
	if maxDate == nil {
		return []quik.SecurityLimit{}, nil
	}
	if err := checkLimitDate(date, *maxDate); err != nil {
		return nil, err
	}
	return s.repo.SelectSecurityLimits(ctx, date)
}

func (s *Service) GetSecurityLimitsOtc(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error) {
	maxDate, err := s.repo.SelectSecurityLimitsOtcMaxDate(ctx)
	if err != nil {
		return nil, err
	}
	if maxDate == nil {
		return []quik.SecurityLimit{}, nil
	}
	if err := checkLimitDate(date, *maxDate); err != nil {
		return nil, err
	}
	return s.repo.SelectSecurityLimitsOtc(ctx, date)
}

func (s *Service) CreateSecurityLimit(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error) {
	maxDate, err := s.repo.SelectSecurityLimitsMaxDate(ctx)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return quik.SecurityLimit{}, err
	}
	if err := checkLimitDate(sec.LoadDate, minRollForwardDate(maxDate)); err != nil {
		return quik.SecurityLimit{}, err
	}

	if sec.SettleCode == "" {
		sec.SettleCode = models.SettleCodeTx
	}

	err = sec.SettleCode.Validate()
	if err != nil {
		return quik.SecurityLimit{}, fmt.Errorf("%w: %s", models.ErrBusinessValidation, err.Error())
	}

	created, err := s.repo.InsertSecurityLimit(ctx, sec)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return quik.SecurityLimit{}, fmt.Errorf("%w: некорректное имя фирмы", models.ErrBusinessValidation)
		}
		if errors.Is(err, models.ErrConflict) {
			return quik.SecurityLimit{}, fmt.Errorf("%w: load_date=%s client_code=%s ticker=%s trade_account=%s settle_code=%s firm_name=%s",
				models.ErrConflict,
				sec.LoadDate.Format(dates.ISODateFormat),
				sec.ClientCode,
				sec.Ticker,
				sec.TradeAccount,
				sec.SettleCode,
				sec.FirmName)
		}
		return quik.SecurityLimit{}, err
	}
	return created, nil
}

func (s *Service) CreateSecurityLimitOtc(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error) {
	maxDate, err := s.repo.SelectSecurityLimitsOtcMaxDate(ctx)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return quik.SecurityLimit{}, err
	}
	if err := checkLimitDate(sec.LoadDate, minRollForwardDate(maxDate)); err != nil {
		return quik.SecurityLimit{}, err
	}

	sec.TradeAccount = "OTC"
	if sec.SettleCode == "" {
		sec.SettleCode = models.SettleCodeTx
	}

	err = sec.SettleCode.Validate()
	if err != nil {
		return quik.SecurityLimit{}, fmt.Errorf("%w: %s", models.ErrBusinessValidation, err.Error())
	}

	created, err := s.repo.InsertSecurityLimitOtc(ctx, sec)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return quik.SecurityLimit{}, fmt.Errorf("%w: некорректное имя фирмы", models.ErrBusinessValidation)
		}
		if errors.Is(err, models.ErrConflict) {
			return quik.SecurityLimit{}, fmt.Errorf("%w: load_date=%s client_code=%s ticker=%s trade_account=%s settle_code=%s firm_name=%s",
				models.ErrConflict,
				sec.LoadDate.Format(dates.ISODateFormat),
				sec.ClientCode,
				sec.Ticker,
				sec.TradeAccount,
				sec.SettleCode,
				sec.FirmName)
		}
		return quik.SecurityLimit{}, err
	}
	return created, nil
}

func (s *Service) DoRollForwardOtc(ctx context.Context) error {
	return doRollForward(ctx,
		s.repo.SelectSecurityLimitsOtcMaxDate,
		s.repo.InsertSecurityLimitsOtcCopy,
		s.repo.DeleteSecurityLimitsOtc,
	)
}

func (s *Service) DoRollForwardSecurityLimits(ctx context.Context) error {
	return doRollForward(ctx,
		s.repo.SelectSecurityLimitsMaxDate,
		s.repo.InsertSecurityLimitsCopy,
		s.repo.DeleteSecurityLimits,
	)
}
