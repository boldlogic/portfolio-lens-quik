package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (s *Service) CreateSecurityLimit(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error) {

	if strings.TrimSpace(sec.SettleCode.String()) == "" {
		sec.SettleCode = quik.SettleCodeTx
	}

	err := sec.SettleCode.Validate()
	if err != nil {
		return quik.SecurityLimit{}, fmt.Errorf("%w: %s", models.ErrBusinessValidation, err.Error())
	}

	created, err := s.repo.InsertSecurityLimit(ctx, sec)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return quik.SecurityLimit{}, fmt.Errorf("%w: некорректный код фирмы %s", models.ErrBusinessValidation, sec.FirmCode)
		}
		if errors.Is(err, models.ErrConflict) {
			return quik.SecurityLimit{}, fmt.Errorf("%w: clientCode=%s secCode=%s tradeAccount=%s settleCode=%s firmCode=%s",
				models.ErrConflict,
				sec.ClientCode,
				sec.SecCode,
				sec.TradeAccount,
				sec.SettleCode,
				sec.FirmCode)
		}
		return quik.SecurityLimit{}, err
	}
	return created, nil
}

func (s *Service) CreateSecurityLimitOtc(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error) {

	sec.TradeAccount = "OTC"
	if strings.TrimSpace(sec.SettleCode.String()) == "" {
		sec.SettleCode = quik.SettleCodeTx
	}

	err := sec.SettleCode.Validate()
	if err != nil {
		return quik.SecurityLimit{}, fmt.Errorf("%w: %s", models.ErrBusinessValidation, err.Error())
	}

	created, err := s.repo.InsertSecurityLimitOtc(ctx, sec)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return quik.SecurityLimit{}, fmt.Errorf("%w: некорректный код фирмы %s", models.ErrBusinessValidation, sec.FirmCode)
		}
		if errors.Is(err, models.ErrConflict) {
			return quik.SecurityLimit{}, fmt.Errorf("%w: clientCode=%s secCode=%s tradeAccount=%s settleCode=%s firmCode=%s",
				models.ErrConflict,
				sec.ClientCode,
				sec.SecCode,
				sec.TradeAccount,
				sec.SettleCode,
				sec.FirmCode)
		}
		return quik.SecurityLimit{}, err
	}
	return created, nil
}
