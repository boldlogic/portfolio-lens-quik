package v1

import (
	"errors"
	"net/http"

	"github.com/boldlogic/packages/transport/httputils"
	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

func (h *Handler) postMoneyLimit(r *http.Request) (any, string, error) {
	ctx := r.Context()
	req, err := httputils.DecodeRequest[moneyLimitReqDTO](r)
	if err != nil {
		if errors.Is(err, httputils.ErrUnsupportedMediaType) || errors.Is(err, httputils.ErrRequestEntityTooLarge) {
			return nil, err.Error(), err
		}
		return nil, err.Error(), models.ErrValidation
	}

	lim, err := req.toQuik()
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

func (req moneyLimitReqDTO) toQuik() (quik.MoneyLimit, error) {

	return quik.MoneyLimit{
		ClientCode:   req.ClientCode,
		CurrencyCode: req.Currency,
		PositionCode: req.PositionCode,
		SettleCode:   quik.SettleCode(req.SettleCode),
		FirmCode:     req.FirmCode,
		Balance:      req.Balance,
	}, nil

}

type moneyLimitReqDTO struct {
	ClientCode   string          `json:"clientCode" validate:"required,min=1,max=12"`
	Currency     string          `json:"currency" validate:"required,min=1,max=3"`
	PositionCode string          `json:"positionCode" validate:"omitempty,min=1,max=4"`
	SettleCode   string          `json:"settleCode" validate:"omitempty,min=0,max=5"`
	FirmCode     string          `json:"firmCode" validate:"required,min=1,max=12"`
	Balance      decimal.Decimal `json:"balance" validate:"required,decimal"`
}

type moneyLimitDTO struct {
	LoadDate     string          `json:"loadDate"`
	SourceDate   string          `json:"sourceDate"`
	ClientCode   string          `json:"clientCode"`
	Currency     string          `json:"currency"`
	PositionCode string          `json:"positionCode"`
	SettleCode   string          `json:"settleCode"`
	FirmCode     string          `json:"firmCode"`
	FirmName     string          `json:"firmName"`
	Balance      decimal.Decimal `json:"balance"`
}

func moneyLimitToDTO(ml quik.MoneyLimit) moneyLimitDTO {
	return moneyLimitDTO{
		LoadDate:     ml.LoadDate.Format(dates.ISODateFormat),
		SourceDate:   ml.SourceDate.Format(dates.ISODateFormat),
		ClientCode:   ml.ClientCode,
		Currency:     ml.CurrencyCode,
		PositionCode: ml.PositionCode,
		SettleCode:   string(ml.SettleCode),
		FirmCode:     ml.FirmCode,
		FirmName:     ml.FirmName,
		Balance:      ml.Balance,
	}
}
