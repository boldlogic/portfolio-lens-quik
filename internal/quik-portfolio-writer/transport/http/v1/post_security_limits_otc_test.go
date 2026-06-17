package v1

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/boldlogic/packages/transport/httputils"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateSecurityLimitOtc(t *testing.T) {
	t.Parallel()

	loadDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)
	sourceDate := time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local)
	internalSvcErr := errors.New("временная_ошибка_хранилища")

	successStub := func(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error) {
		sec.LoadDate = loadDate
		sec.SourceDate = sourceDate
		sec.TradeAccount = "OTC"
		sec.FirmName = "Фирма брокера"
		return sec, nil
	}

	tests := []struct {
		name       string
		req        *http.Request
		svc        svc
		wantBody   any
		wantDetail string
		wantErr    error
	}{
		{
			name: "успешный_запрос_min_полей",
			req: reqJSON(`{
				"clientCode":"AB12CD",
				"secCode":"OTC_BOND",
				"settleCode":"T0",
				"firmCode": "NC0058900000",
				"balance":2,
				"acquisitionCurrencyCode":"USD"
			}`),
			svc: svc{createSecurityLimitOtc: successStub},
			wantBody: securityLimitDTO{
				LoadDate:                "2025-01-01",
				SourceDate:              "2025-01-02",
				ClientCode:              "AB12CD",
				SecCode:                 "OTC_BOND",
				TradeAccount:            "OTC",
				SettleCode:              "T0",
				FirmCode:                "NC0058900000",
				FirmName:                "Фирма брокера",
				Balance:                 decimal.NewFromFloat(2),
				AcquisitionCurrencyCode: "USD",
				ISIN:                    "",
			},
		},
		{
			name: "успех_с_isin",
			req: reqJSON(`{
				"clientCode":"AB12CD",
				"secCode":"OTC_BOND",
				"settleCode":"T0",
				"firmCode": "NC0058900000",
				"balance":2,
				"acquisitionCurrencyCode":"USD",
				"isin":"RU000A0JX0J2"
			}`),
			svc: svc{createSecurityLimitOtc: successStub},
			wantBody: securityLimitDTO{
				LoadDate:                "2025-01-01",
				SourceDate:              "2025-01-02",
				ClientCode:              "AB12CD",
				SecCode:                 "OTC_BOND",
				TradeAccount:            "OTC",
				SettleCode:              "T0",
				FirmCode:                "NC0058900000",
				FirmName:                "Фирма брокера",
				Balance:                 decimal.NewFromFloat(2),
				AcquisitionCurrencyCode: "USD",
				ISIN:                    "RU000A0JX0J2",
			},
		},
		{
			name: "конфликт_ключа",
			req: reqJSON(`{
				"clientCode":"AB12CD",
				"secCode":"OTC_BOND",
				"firmCode": "NC0058900000",
				"balance":2,
				"acquisitionCurrencyCode":"USD"
			}`),
			svc:        svc{err: models.ErrConflict},
			wantBody:   nil,
			wantDetail: models.ErrConflict.Error(),
			wantErr:    models.ErrConflict,
		},
		{
			name: "бизнес_валидация_ErrBusinessValidation",
			req: reqJSON(`{
				"clientCode":"AB12CD",
				"secCode":"OTC_BOND",
				"firmCode": "NC0058900000",
				"balance":2,
				"acquisitionCurrencyCode":"USD"
			}`),
			svc:        svc{err: models.ErrBusinessValidation},
			wantBody:   nil,
			wantDetail: models.ErrBusinessValidation.Error(),
			wantErr:    models.ErrBusinessValidation,
		},
		{
			name: "внутренняя_ошибка_сервиса",
			req: reqJSON(`{
				"clientCode":"AB12CD",
				"secCode":"OTC_BOND",
				"firmCode": "NC0058900000",
				"balance":1,
				"acquisitionCurrencyCode":"USD"
			}`),
			svc:        svc{err: internalSvcErr},
			wantBody:   nil,
			wantDetail: "",
			wantErr:    internalSvcErr,
		},
		{
			name: "UnsupportedMediaType",
			req: httptest.NewRequest(http.MethodPost, exampleURL, bytes.NewBufferString(
				`{"clientCode":"AB12CD","secCode":"OTC_BOND","firmCode":"NC0058900000","balance":1,"acquisitionCurrencyCode":"USD"}`,
			)),
			wantBody:   nil,
			wantDetail: "Content-Type",
			wantErr:    httputils.ErrUnsupportedMediaType,
		},
		{
			name:       "битый_json",
			req:        reqJSON(`{"clientCode":"AB12CD"`),
			wantBody:   nil,
			wantDetail: "unexpected EOF",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "пустой_clientCode",
			req:        reqJSON(`{"clientCode":"","secCode":"OTC_BOND","firmCode":"NC0058900000","balance":1,"acquisitionCurrencyCode":"USD"}`),
			wantBody:   nil,
			wantDetail: "clientCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "clientCode_длина_>12",
			req:        reqJSON(`{"clientCode":"1234567890123","secCode":"OTC_BOND","firmCode":"NC0058900000","balance":1,"acquisitionCurrencyCode":"USD"}`),
			wantBody:   nil,
			wantDetail: "clientCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "пустой_ticker",
			req:        reqJSON(`{"clientCode":"AB12CD","secCode":"","firmCode":"NC0058900000","balance":1,"acquisitionCurrencyCode":"USD"}`),
			wantBody:   nil,
			wantDetail: "secCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "пустой_firmCode",
			req:        reqJSON(`{"clientCode":"AB12CD","secCode":"OTC_BOND","firmCode":"","balance":1,"acquisitionCurrencyCode":"USD"}`),
			wantBody:   nil,
			wantDetail: "firmCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "лишнее_поле_tradeAccount_DecodeJSONStrict",
			req:        reqJSON(`{"clientCode":"AB12CD","secCode":"OTC_BOND","firmCode":"NC0058900000","balance":1,"acquisitionCurrencyCode":"USD","tradeAccount":"L01-00000F00"}`),
			wantBody:   nil,
			wantDetail: "некорректный формат JSON",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "ticker_длина_>12",
			req:        reqJSON(`{"clientCode":"AB12CD","secCode":"SBERLONGTICKER","firmCode":"NC0058900000","balance":1,"acquisitionCurrencyCode":"USD"}`),
			wantBody:   nil,
			wantDetail: "secCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "firmCode_длина_>12",
			req:        reqJSON(`{"clientCode":"AB12CD","secCode":"OTC_BOND","firmCode":"NC00589000001","balance":1,"acquisitionCurrencyCode":"USD"}`),
			wantBody:   nil,
			wantDetail: "firmCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "settleCode_длина_>5",
			req:        reqJSON(`{"clientCode":"AB12CD","secCode":"OTC_BOND","settleCode":"123456","firmCode":"NC0058900000","balance":1,"acquisitionCurrencyCode":"USD"}`),
			wantBody:   nil,
			wantDetail: "settleCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "acquisitionCcy_длина_>3",
			req:        reqJSON(`{"clientCode":"AB12CD","secCode":"OTC_BOND","firmCode":"NC0058900000","balance":1,"acquisitionCurrencyCode":"USDD"}`),
			wantBody:   nil,
			wantDetail: "acquisitionCurrencyCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "isin_длина_>12",
			req:        reqJSON(`{"clientCode":"AB12CD","secCode":"OTC_BOND","firmCode":"NC0058900000","balance":1,"acquisitionCurrencyCode":"USD","isin":"RU000A0JX0J2X"}`),
			wantBody:   nil,
			wantDetail: "isin",
			wantErr:    models.ErrValidation,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := newTestHandler(tt.svc)
			body, detail, err := h.CreateSecurityLimitOtc(tt.req)
			assert.Equal(t, tt.wantBody, body)
			assert.Contains(t, detail, tt.wantDetail)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
