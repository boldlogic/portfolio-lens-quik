package v1

import (
	"errors"
	"net/http"

	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/httputils"
	"github.com/shopspring/decimal"
)

func (h *Handler) CreateSecurityLimitOtc(r *http.Request) (any, string, error) {
	ctx := r.Context()
	req, err := httputils.DecodeAndValidate[securityLimitOtcReqDTO](r)
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

	created, err := h.service.CreateSecurityLimitOtc(ctx, lim)
	if err != nil {
		if errors.Is(err, models.ErrBusinessValidation) || errors.Is(err, models.ErrConflict) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return securityLimitToDTO(created), "", nil
}

func (req securityLimitOtcReqDTO) convertToSecurityLimit() (quik.SecurityLimit, error) {
	date, err := dates.ParseWithDefaultNow(req.LoadDate, dates.ISODateFormat)
	if err != nil {
		return quik.SecurityLimit{}, err
	}
	var isin *string
	if req.ISIN != "" {
		isin = &req.ISIN
	}
	return quik.SecurityLimit{
		LoadDate:       date,
		ClientCode:     req.ClientCode,
		Ticker:         req.Ticker,
		SettleCode:     models.SettleCode(req.SettleCode),
		FirmName:       req.FirmName,
		Balance:        decimal.NewFromFloat(req.Balance),
		AcquisitionCcy: req.AcquisitionCcy,
		ISIN:           isin,
	}, nil
}

type securityLimitOtcReqDTO struct {
	LoadDate       string  `json:"loadDate" validate:"omitempty"`
	ClientCode     string  `json:"clientCode" validate:"required,min=1,max=12"`
	Ticker         string  `json:"ticker" validate:"required,min=1,max=12"`
	SettleCode     string  `json:"settleCode" validate:"omitempty,min=0,max=5"`
	FirmName       string  `json:"firmName" validate:"required,min=1,max=128"`
	Balance        float64 `json:"balance" validate:"min=-999999999999,max=999999999999"`
	AcquisitionCcy string  `json:"acquisitionCcy" validate:"omitempty,min=1,max=3"`
	ISIN           string  `json:"isin" validate:"omitempty,min=1,max=12"`
}
