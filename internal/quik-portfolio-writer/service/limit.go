package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/boldlogic/portfolio-lens-currency/pkg/currencies"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (s *Service) CreateLimit(ctx context.Context, l quik.Limit) (quik.Limit, error) {
	var err error
	switch l.Type {
	case quik.LimitTypeMoney:
		err = prepareMoneyLimit(&l)
	case quik.LimitTypeSecurities:
		err = prepareSecurityLimit(&l, false)
	case quik.LimitTypeSecuritiesOtc:
		err = prepareSecurityLimit(&l, true)
	default:
		return quik.Limit{}, fmt.Errorf("%w: неподдерживаемый тип лимита", models.ErrBusinessValidation)
	}
	if err != nil {
		return quik.Limit{}, err
	}

	created, err := s.repo.InsertLimit(ctx, l)
	if err != nil {
		return quik.Limit{}, mapCreateLimitError(l, err)
	}
	return created, nil
}

func prepareMoneyLimit(l *quik.Limit) error {
	if l.CurrencyCode == nil || strings.TrimSpace(*l.CurrencyCode) == "" {
		return fmt.Errorf("%w: не указана валюта", models.ErrBusinessValidation)
	}

	normCcy, err := currencies.ParseCurrencyCode(*l.CurrencyCode)
	if err != nil {
		return fmt.Errorf("%w: %s", models.ErrBusinessValidation, err.Error())
	}
	norm := normCcy.String()
	l.CurrencyCode = &norm

	if strings.TrimSpace(l.SettleCode.String()) == "" {
		l.SettleCode = quik.SettleCodeTx
	}
	if err := l.SettleCode.Validate(); err != nil {
		return fmt.Errorf("%w: %s", models.ErrBusinessValidation, err.Error())
	}
	return nil
}

func prepareSecurityLimit(l *quik.Limit, otc bool) error {
	if otc {
		tradeAccount := "OTC"
		l.TradeAccount = &tradeAccount
	}

	if strings.TrimSpace(l.SettleCode.String()) == "" {
		l.SettleCode = quik.SettleCodeTx
	}
	if err := l.SettleCode.Validate(); err != nil {
		return fmt.Errorf("%w: %s", models.ErrBusinessValidation, err.Error())
	}
	return nil
}

func mapCreateLimitError(l quik.Limit, err error) error {
	if errors.Is(err, models.ErrNotFound) {
		return fmt.Errorf("%w: некорректный код фирмы %s", models.ErrBusinessValidation, l.FirmCode)
	}
	if !errors.Is(err, models.ErrConflict) {
		return err
	}

	switch l.Type {
	case quik.LimitTypeMoney:
		positionCode := ""
		if l.PositionCode != nil {
			positionCode = *l.PositionCode
		}
		currency := ""
		if l.CurrencyCode != nil {
			currency = *l.CurrencyCode
		}
		return fmt.Errorf("%w: clientCode=%s currency=%s positionCode=%s settleCode=%s firmCode=%s",
			models.ErrConflict,
			l.ClientCode,
			currency,
			positionCode,
			l.SettleCode,
			l.FirmCode)
	default:
		secCode := ""
		if l.SecCode != nil {
			secCode = *l.SecCode
		}
		tradeAccount := ""
		if l.TradeAccount != nil {
			tradeAccount = *l.TradeAccount
		}
		return fmt.Errorf("%w: clientCode=%s secCode=%s tradeAccount=%s settleCode=%s firmCode=%s",
			models.ErrConflict,
			l.ClientCode,
			secCode,
			tradeAccount,
			l.SettleCode,
			l.FirmCode)
	}
}
