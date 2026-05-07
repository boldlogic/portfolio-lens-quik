package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/boldlogic/packages/utils/dates"
	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
)

func (h *Handler) extractDateQueryParam(r *http.Request) (time.Time, error) {
	dateReq := r.URL.Query().Get("date")
	return dates.ParseWithDefaultNow(dateReq, dates.ISODateFormat)
}

func (h *Handler) GetLimits(r *http.Request) (any, string, error) {
	ctx := r.Context()
	date, err := h.extractDateQueryParam(r)
	if err != nil {
		return nil, err.Error(), md.ErrValidation
	}
	lim, err := h.service.GetLimits(ctx, date)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, err.Error(), err
		}

		return nil, "", err
	}

	return limitsToResp(lim), "", nil
}

func limitsToResp(limits []quik.Limit) []limitDTO {
	if len(limits) == 0 {
		return []limitDTO{}
	}
	resp := make([]limitDTO, 0, len(limits))
	for _, l := range limits {
		resp = append(resp, limitToDTO(l))
	}
	return resp
}

func limitToDTO(limit quik.Limit) limitDTO {
	var isin string
	if limit.ISIN != nil {
		isin = *limit.ISIN
	}
	return limitDTO{
		LimitType:      string(limit.LimitType),
		LoadDate:       limit.LoadDate.Format(dates.ISODateFormat),
		SourceDate:     limit.SourceDate.Format(dates.ISODateFormat),
		ClientCode:     limit.ClientCode,
		Instrument:     limit.InstrumentCode,
		ISIN:           isin,
		SettleCode:     string(limit.SettleCode),
		FirmCode:       limit.FirmCode,
		FirmName:       limit.FirmName,
		Balance:        limit.Balance,
		AcquisitionCcy: limit.AcquisitionCcy,
	}
}

type limitDTO struct {
	LimitType      string          `json:"limitType"`
	LoadDate       string          `json:"loadDate"`
	SourceDate     string          `json:"sourceDate"`
	ClientCode     string          `json:"clientCode"`
	Instrument     string          `json:"instrument"`
	ISIN           string          `json:"isin,omitempty"`
	SettleCode     string          `json:"settleCode"`
	FirmCode       string          `json:"firmCode"`
	FirmName       string          `json:"firmName"`
	Balance        decimal.Decimal `json:"balance"`
	AcquisitionCcy string          `json:"acquisitionCcy"`
}
