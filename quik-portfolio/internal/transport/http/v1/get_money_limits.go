package v1

import (
	"errors"
	"net/http"

	"github.com/boldlogic/packages/utils/dates"
	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

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

func moneyLimitsToResp(mls []quik.MoneyLimit) []moneyLimitDTO {
	if len(mls) == 0 {
		return []moneyLimitDTO{}
	}

	resp := make([]moneyLimitDTO, 0, len(mls))
	for _, ml := range mls {
		resp = append(resp, moneyLimitToDTO(ml))
	}
	return resp
}

func moneyLimitToDTO(ml quik.MoneyLimit) moneyLimitDTO {
	return moneyLimitDTO{
		LoadDate:     ml.LoadDate.Format(dates.ISODateFormat),
		SourceDate:   ml.SourceDate.Format(dates.ISODateFormat),
		ClientCode:   ml.ClientCode,
		Currency:     ml.Currency,
		PositionCode: ml.PositionCode,
		SettleCode:   string(ml.SettleCode),
		FirmCode:     ml.FirmCode,
		FirmName:     ml.FirmName,
		Balance:      ml.Balance,
	}
}

func (h *Handler) GetMoneyLimits(r *http.Request) (any, string, error) {
	ctx := r.Context()
	date, err := h.extractDateQueryParam(r)
	if err != nil {
		return nil, err.Error(), md.ErrValidation
	}

	mls, err := h.service.GetMoneyLimits(ctx, date)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}

		return nil, "", err
	}

	return moneyLimitsToResp(mls), "", nil
}
