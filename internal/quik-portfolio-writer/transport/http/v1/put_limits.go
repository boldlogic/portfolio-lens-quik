package v1

import (
	"errors"
	"net/http"

	"github.com/boldlogic/packages/transport/httputils"
	intmodels "github.com/boldlogic/portfolio-lens-quik/internal/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type limitReq struct {
	Type                    string          `json:"limitType" validate:"oneof=money security security_otc"`
	ClientCode              string          `json:"clientCode" validate:"required,min=1,max=12"`
	Ticker                  string          `json:"ticker" validate:"required,min=1,max=12"`
	PositionCode            *string         `json:"positionCode" validate:"omitempty,min=1,max=4"`
	SettleCode              string          `json:"settleCode" validate:"omitempty,oneof=T0 T1 T2 Tx"`
	TradeAccount            *string         `json:"tradeAccount" validate:"omitempty,min=1,max=12"`
	FirmCode                string          `json:"firmCode" validate:"required,min=1,max=12"`
	Balance                 decimal.Decimal `json:"balance" validate:"required"`
	AcquisitionCurrencyCode *string         `json:"acquisitionCurrencyCode" validate:"omitempty,min=1,max=3"`
	ISIN                    *string         `json:"isin" validate:"omitempty,min=1,max=12"`
}

func (h *Handler) createLimit(r *http.Request) (any, string, error) {
	ctx := r.Context()

	req, err := httputils.DecodeRequest[limitReq](r)
	if err != nil {
		if errors.Is(err, httputils.ErrUnsupportedMediaType) || errors.Is(err, httputils.ErrRequestEntityTooLarge) {
			return nil, err.Error(), err
		}
		return nil, err.Error(), models.ErrValidation
	}

	lim := intmodels.LimitInput{
		Type:                    req.Type,
		ClientCode:              req.ClientCode,
		Ticker:                  req.Ticker,
		SettleCode:              req.SettleCode,
		PositionCode:            req.PositionCode,
		TradeAccount:            req.TradeAccount,
		FirmCode:                req.FirmCode,
		Balance:                 req.Balance,
		AcquisitionCurrencyCode: req.AcquisitionCurrencyCode,
		ISIN:                    req.ISIN,
	}

	err = h.service.UpsertLimit(ctx, lim)
	if err != nil {
		h.logger.Error("", zap.Error(err))

		if errors.Is(err, models.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}

	return nil, "", nil
}
