package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/boldlogic/portfolio-lens-currency/pkg/currencies"
	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

// CreateMoneyLimit
func (s *Service) CreateMoneyLimit(ctx context.Context, ml quik.MoneyLimit) (quik.MoneyLimit, error) {

	normCcy, err := currencies.ParseCurrencyCode(ml.Currency)
	if err != nil {
		return quik.MoneyLimit{}, fmt.Errorf("%w: %s", md.ErrBusinessValidation, err.Error())
	}
	ml.Currency = normCcy.String()

	if strings.TrimSpace(ml.SettleCode.String()) == "" {
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
			return quik.MoneyLimit{}, fmt.Errorf("%w: clientCode=%s currency=%s positionCode=%s settleCode=%s firmCode=%s",
				md.ErrConflict,
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
