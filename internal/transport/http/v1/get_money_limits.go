package v1

import (
	"errors"
	"net/http"

	"github.com/boldlogic/quik-portfolio/internal/models"
	md "github.com/boldlogic/quik-portfolio/pkg/models"
)

func (h *Handler) GetMoneyLimits(r *http.Request) (any, string, error) {
	ctx := r.Context()
	date, err := h.readGetLimitsRequest(r)
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

func moneyLimitsToResp(mls []models.MoneyLimit) []moneyLimitDTO {
	if len(mls) == 0 {
		return []moneyLimitDTO{}
	}

	resp := make([]moneyLimitDTO, 0, len(mls))
	for _, ml := range mls {
		resp = append(resp, moneyLimitToDTO(ml))
	}
	return resp
}

func moneyLimitToDTO(ml models.MoneyLimit) moneyLimitDTO {
	return moneyLimitDTO{
		LoadDate:     ml.LoadDate.Format(md.ISODateFormat),
		SourceDate:   ml.SourceDate.Format(md.ISODateFormat),
		ClientCode:   ml.ClientCode,
		Currency:     ml.Currency,
		PositionCode: ml.PositionCode,
		SettleCode:   string(ml.SettleCode),
		FirmCode:     ml.FirmCode,
		FirmName:     ml.FirmName,
		Balance:      ml.Balance.InexactFloat64(),
	}
}

type moneyLimitDTO struct {
	LoadDate     string  `json:"loadDate"`
	SourceDate   string  `json:"sourceDate"`
	ClientCode   string  `json:"clientCode"`
	Currency     string  `json:"currency"`
	PositionCode string  `json:"positionCode"`
	SettleCode   string  `json:"settleCode"`
	FirmCode     string  `json:"firmCode"`
	FirmName     string  `json:"firmName"`
	Balance      float64 `json:"balance"`
}
