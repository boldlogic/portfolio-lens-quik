package v1

import (
	"errors"
	"net/http"

	"github.com/boldlogic/packages/transport/httputils"
	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type limitResp struct {
	Type                    string          `json:"limitType"`
	LoadDate                string          `json:"loadDate"`
	SourceDate              string          `json:"sourceDate"`
	ClientCode              string          `json:"clientCode"`
	CurrencyCode            *string         `json:"currency,omitempty"`
	SecCode                 *string         `json:"secCode,omitempty"`
	PositionCode            *string         `json:"positionCode,omitempty"`
	SettleCode              string          `json:"settleCode"`
	TradeAccount            *string         `json:"tradeAccount,omitempty"`
	FirmCode                string          `json:"firmCode"`
	FirmName                string          `json:"firmName"`
	Balance                 decimal.Decimal `json:"balance"`
	AcquisitionCurrencyCode *string         `json:"acquisitionCurrencyCode,omitempty"`
	ISIN                    *string         `json:"isin,omitempty"`
	ShortName               *string         `json:"shortName,omitempty"`
}

func limitToResp(l quik.Limit) limitResp {
	return limitResp{
		Type:                    string(l.Type),
		LoadDate:                l.LoadDate.Format(dates.ISODateFormat),
		SourceDate:              l.SourceDate.Format(dates.ISODateFormat),
		ClientCode:              l.ClientCode,
		CurrencyCode:            l.CurrencyCode,
		SecCode:                 l.SecCode,
		PositionCode:            l.PositionCode,
		SettleCode:              string(l.SettleCode),
		TradeAccount:            l.TradeAccount,
		FirmCode:                l.FirmCode,
		FirmName:                l.FirmName,
		Balance:                 l.Balance,
		AcquisitionCurrencyCode: l.AcquisitionCurrencyCode,
		ISIN:                    l.ISIN,
		ShortName:               l.ShortName,
	}
}

type limitReq struct {
	Type                    string          `json:"limitType" validate:"oneof=money security security_otc"`
	ClientCode              string          `json:"clientCode" validate:"required,min=1,max=12"`
	CurrencyCode            *string         `json:"currency" validate:"omitempty,min=1,max=3"`
	SecCode                 *string         `json:"secCode" validate:"omitempty,min=1,max=12"`
	PositionCode            *string         `json:"positionCode,omitempty,min=1,max=4"`
	SettleCode              string          `json:"settleCode" validate:"omitempty,min=0,max=5"`
	TradeAccount            *string         `json:"tradeAccount" validate:"omitempty,min=1,max=12"`
	FirmCode                string          `json:"firmCode" validate:"required,min=1,max=12"`
	Balance                 decimal.Decimal `json:"balance" validate:"required"`
	AcquisitionCurrencyCode *string         `json:"acquisitionCurrencyCode" validate:"omitempty,min=1,max=3"`
	ISIN                    *string         `json:"isin" validate:"omitempty,min=1,max=12"`
}

func (r limitReq) convert() (quik.Limit, error) {
	lt := quik.LimitType(r.Type)

	settleCode := quik.SettleCode(r.SettleCode)
	if r.SettleCode != "" {
		if err := settleCode.Validate(); err != nil {
			return quik.Limit{}, err
		}
	}

	return quik.Limit{
		Type:                    lt,
		ClientCode:              r.ClientCode,
		CurrencyCode:            r.CurrencyCode,
		SecCode:                 r.SecCode,
		PositionCode:            r.PositionCode,
		SettleCode:              settleCode,
		TradeAccount:            r.TradeAccount,
		FirmCode:                r.FirmCode,
		Balance:                 r.Balance,
		AcquisitionCurrencyCode: r.AcquisitionCurrencyCode,
		ISIN:                    r.ISIN,
	}, nil
}

func (h *Handler) createLimit(r *http.Request) (any, string, error) {
	ctx := r.Context()

	req, err := httputils.DecodeRequest[limitReq](r)
	if err != nil {
		if errors.Is(err, httputils.ErrUnsupportedMediaType) || errors.Is(err, httputils.ErrRequestEntityTooLarge) {
			return nil, err.Error(), err
		}
		h.logger.Error("", zap.Error(err))
		return nil, err.Error(), models.ErrValidation
	}
	lim, err := req.convert()
	if err != nil {
		h.logger.Error("", zap.Error(err))

		return nil, err.Error(), models.ErrValidation
	}
	created, err := h.service.CreateLimit(ctx, lim)
	if err != nil {
		h.logger.Error("", zap.Error(err))

		if errors.Is(err, models.ErrBusinessValidation) || errors.Is(err, models.ErrConflict) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	created.Type = lim.Type
	return limitToResp(created), "", nil
}
