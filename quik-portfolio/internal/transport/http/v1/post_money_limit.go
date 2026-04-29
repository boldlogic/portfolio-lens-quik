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

func (h *Handler) CreateMoneyLimit(r *http.Request) (any, string, error) {
	ctx := r.Context()
	req, err := httputils.DecodeAndValidate[moneyLimitReqDTO](r)
	if err != nil {
		if errors.Is(err, httputils.ErrUnsupportedMediaType) || errors.Is(err, httputils.ErrRequestEntityTooLarge) {
			return nil, err.Error(), err
		}
		return nil, err.Error(), models.ErrValidation
	}

	lim, err := req.convertToMoneyLimit()
	if err != nil {
		return nil, err.Error(), models.ErrValidation
	}

	created, err := h.service.CreateMoneyLimit(ctx, lim)
	if err != nil {
		if errors.Is(err, models.ErrBusinessValidation) || errors.Is(err, models.ErrConflict) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return moneyLimitToDTO(created), "", nil
}

func (req moneyLimitReqDTO) convertToMoneyLimit() (quik.MoneyLimit, error) {
	date, err := dates.ParseWithDefaultNow(req.LoadDate, dates.ISODateFormat)
	if err != nil {
		return quik.MoneyLimit{}, err
	}

	return quik.MoneyLimit{
		LoadDate:     date,
		ClientCode:   req.ClientCode,
		Currency:     req.Currency,
		PositionCode: req.PositionCode,
		SettleCode:   models.SettleCode(req.SettleCode),
		FirmName:     req.FirmName,
		Balance:      decimal.NewFromFloat(req.Balance),
	}, nil

}

type moneyLimitReqDTO struct {
	LoadDate     string  `json:"loadDate" validate:"omitempty"`
	ClientCode   string  `json:"clientCode" validate:"required,min=1,max=12"`
	Currency     string  `json:"currency" validate:"required,min=1,max=3"`
	PositionCode string  `json:"positionCode" validate:"omitempty,min=1,max=4"`
	SettleCode   string  `json:"settleCode" validate:"omitempty,min=0,max=5"`
	FirmName     string  `json:"firmName" validate:"required,min=1,max=128"`
	Balance      float64 `json:"balance" validate:"min=-999999999999,max=999999999999"`
}
