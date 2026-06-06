package v1

import (
	"errors"
	"net/http"

	"github.com/boldlogic/packages/transport/httputils"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

func (h *Handler) CreateSecurityLimit(r *http.Request) (any, string, error) {
	ctx := r.Context()
	req, err := httputils.DecodeRequest[securityLimitReqDTO](r)
	if err != nil {
		if errors.Is(err, httputils.ErrUnsupportedMediaType) || errors.Is(err, httputils.ErrRequestEntityTooLarge) {
			return nil, err.Error(), err
		}
		return nil, err.Error(), models.ErrValidation
	}

	lim, err := req.convertToSecurityLimit()
	if err != nil {
		return nil, err.Error(), models.ErrValidation
	}

	created, err := h.service.CreateSecurityLimit(ctx, lim)
	if err != nil {
		if errors.Is(err, models.ErrBusinessValidation) || errors.Is(err, models.ErrConflict) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return securityLimitToDTO(created), "", nil
}

func (req securityLimitReqDTO) convertToSecurityLimit() (quik.SecurityLimit, error) {

	return quik.SecurityLimit{
		ClientCode:     req.ClientCode,
		Ticker:         req.Ticker,
		TradeAccount:   req.TradeAccount,
		SettleCode:     quik.SettleCode(req.SettleCode),
		FirmCode:       req.FirmCode,
		Balance:        req.Balance,
		AcquisitionCcy: req.AcquisitionCcy,
		ISIN:           req.ISIN,
	}, nil
}

type securityLimitReqDTO struct {
	ClientCode     string          `json:"clientCode" validate:"required,min=1,max=12"`
	Ticker         string          `json:"ticker" validate:"required,min=1,max=12"`
	TradeAccount   string          `json:"tradeAccount" validate:"required,min=1,max=12"`
	SettleCode     string          `json:"settleCode" validate:"omitempty,min=0,max=5"`
	FirmCode       string          `json:"firmCode" validate:"required,min=1,max=12"`
	Balance        decimal.Decimal `json:"balance"`
	AcquisitionCcy string          `json:"acquisitionCcy" validate:"omitempty,min=1,max=3"`
	ISIN           string          `json:"isin" validate:"omitempty,min=1,max=12"`
}
