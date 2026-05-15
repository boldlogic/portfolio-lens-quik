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
	"github.com/boldlogic/packages/utils/dates"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateSecurityLimit(t *testing.T) {
	t.Parallel()

	sourceDate := time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local)
	internalSvcErr := errors.New("временная_ошибка_хранилища")

	successStub := func(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error) {
		sec.SourceDate = sourceDate
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
			name: "успешный_запрос",
			req: reqJSON(`{
				"loadDate":"2025-01-01",
				"clientCode":"AB12CD",
				"ticker":"SBER",
				"tradeAccount":"L01-00000F00",
				"settleCode":"T2",
				"firmCode": "NC0058900000",
				"balance":10.5,
				"acquisitionCcy":"RUB",
				"isin":"RU000A0JX0J2"
			}`),
			svc: svc{createSecurityLimit: successStub},
			wantBody: securityLimitDTO{
				LoadDate:       "2025-01-01",
				SourceDate:     "2025-01-02",
				ClientCode:     "AB12CD",
				Ticker:         "SBER",
				TradeAccount:   "L01-00000F00",
				SettleCode:     "T2",
				FirmCode:       "NC0058900000",
				FirmName:       "Фирма брокера",
				Balance:        decimal.NewFromFloat(10.5),
				AcquisitionCcy: "RUB",
				ISIN:           "RU000A0JX0J2",
			},
		},
		{
			name: "успешный_запрос_без_isin",
			req: reqJSON(`{
				"loadDate":"2025-01-01",
				"clientCode":"AB12CD",
				"ticker":"SBER",
				"tradeAccount":"L01-00000F00",
				"settleCode":"T2",
				"firmCode": "NC0058900000",
				"balance":10.5,
				"acquisitionCcy":"RUB"
			}`),
			svc: svc{createSecurityLimit: successStub},
			wantBody: securityLimitDTO{
				LoadDate:       "2025-01-01",
				SourceDate:     "2025-01-02",
				ClientCode:     "AB12CD",
				Ticker:         "SBER",
				TradeAccount:   "L01-00000F00",
				SettleCode:     "T2",
				FirmCode:       "NC0058900000",
				FirmName:       "Фирма брокера",
				Balance:        decimal.NewFromFloat(10.5),
				AcquisitionCcy: "RUB",
				ISIN:           "",
			},
		},
		{
			name: "конфликт_ключа",
			req: reqJSON(`{
				"loadDate":"2025-01-01",
				"clientCode":"AB12CD",
				"ticker":"SBER",
				"tradeAccount":"L01-00000F00",
				"firmCode": "NC0058900000",
				"balance":10.5
			}`),
			svc:        svc{err: models.ErrConflict},
			wantBody:   nil,
			wantDetail: models.ErrConflict.Error(),
			wantErr:    models.ErrConflict,
		},
		{
			name: "бизнес_валидация_ErrBusinessValidation",
			req: reqJSON(`{
				"loadDate":"2025-01-01",
				"clientCode":"AB12CD",
				"ticker":"SBER",
				"tradeAccount":"L01-00000F00",
				"firmCode": "NC0058900000",
				"balance":10.5
			}`),
			svc:        svc{err: models.ErrBusinessValidation},
			wantBody:   nil,
			wantDetail: models.ErrBusinessValidation.Error(),
			wantErr:    models.ErrBusinessValidation,
		},
		{
			name: "внутренняя_ошибка_сервиса",
			req: reqJSON(`{
				"loadDate":"2025-01-01",
				"clientCode":"AB12CD",
				"ticker":"SBER",
				"tradeAccount":"L01-00000F00",
				"firmCode": "NC0058900000",
				"balance":1
			}`),
			svc:        svc{err: internalSvcErr},
			wantBody:   nil,
			wantDetail: "",
			wantErr:    internalSvcErr,
		},
		{
			name: "UnsupportedMediaType",
			req: httptest.NewRequest(http.MethodPost, exampleURL, bytes.NewBufferString(
				`{"loadDate":"2025-01-01","clientCode":"AB12CD","ticker":"SBER","tradeAccount":"L01-00000F00","firmCode":"NC0058900000","balance":1}`,
			)),
			wantBody:   nil,
			wantDetail: "Content-Type",
			wantErr:    httputils.ErrUnsupportedMediaType,
		},
		{
			name:       "битый_json",
			req:        reqJSON(`{"loadDate":"2025-01-01"`),
			wantBody:   nil,
			wantDetail: "unexpected EOF",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "некорректная_дата",
			req:        reqJSON(`{"loadDate":"2025-01","clientCode":"AB12CD","ticker":"SBER","tradeAccount":"L01-00000F00","firmCode":"NC0058900000","balance":1}`),
			wantBody:   nil,
			wantDetail: dates.ErrWrongDateFormat.Error(),
			wantErr:    models.ErrValidation,
		},
		{
			name:       "пустой_tradeAccount",
			req:        reqJSON(`{"loadDate":"2025-01-01","clientCode":"AB12CD","ticker":"SBER","tradeAccount":"","firmCode":"NC0058900000","balance":1}`),
			wantBody:   nil,
			wantDetail: "tradeAccount",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "tradeAccount_длина_>12",
			req:        reqJSON(`{"loadDate":"2025-01-01","clientCode":"AB12CD","ticker":"SBER","tradeAccount":"L01-00000F000","firmCode":"NC0058900000","balance":1}`),
			wantBody:   nil,
			wantDetail: "tradeAccount",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "пустой_clientCode",
			req:        reqJSON(`{"loadDate":"2025-01-01","clientCode":"","ticker":"SBER","tradeAccount":"L01-00000F00","firmCode":"NC0058900000","balance":1}`),
			wantBody:   nil,
			wantDetail: "clientCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "clientCode_длина_>12",
			req:        reqJSON(`{"loadDate":"2025-01-01","clientCode":"1234567890123","ticker":"SBER","tradeAccount":"L01-00000F00","firmCode":"NC0058900000","balance":1}`),
			wantBody:   nil,
			wantDetail: "clientCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "пустой_ticker",
			req:        reqJSON(`{"loadDate":"2025-01-01","clientCode":"AB12CD","ticker":"","tradeAccount":"L01-00000F00","firmCode":"NC0058900000","balance":1}`),
			wantBody:   nil,
			wantDetail: "ticker",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "пустой_firmCode",
			req:        reqJSON(`{"loadDate":"2025-01-01","clientCode":"AB12CD","ticker":"SBER","tradeAccount":"L01-00000F00","firmCode":"","balance":1}`),
			wantBody:   nil,
			wantDetail: "firmCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "ticker_длина_>12",
			req:        reqJSON(`{"loadDate":"2025-01-01","clientCode":"AB12CD","ticker":"SBERLONGTICKER","tradeAccount":"L01-00000F00","firmCode":"NC0058900000","balance":1}`),
			wantBody:   nil,
			wantDetail: "ticker",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "firmCode_длина_>12",
			req:        reqJSON(`{"loadDate":"2025-01-01","clientCode":"AB12CD","ticker":"SBER","tradeAccount":"L01-00000F00","firmCode":"NC00589000001","balance":1}`),
			wantBody:   nil,
			wantDetail: "firmCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "settleCode_длина_>5",
			req:        reqJSON(`{"loadDate":"2025-01-01","clientCode":"AB12CD","ticker":"SBER","tradeAccount":"L01-00000F00","settleCode":"123456","firmCode":"NC0058900000","balance":1}`),
			wantBody:   nil,
			wantDetail: "settleCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "acquisitionCcy_длина_>3",
			req:        reqJSON(`{"loadDate":"2025-01-01","clientCode":"AB12CD","ticker":"SBER","tradeAccount":"L01-00000F00","firmCode":"NC0058900000","balance":1,"acquisitionCcy":"USDD"}`),
			wantBody:   nil,
			wantDetail: "acquisitionCcy",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "isin_длина_>12",
			req:        reqJSON(`{"loadDate":"2025-01-01","clientCode":"AB12CD","ticker":"SBER","tradeAccount":"L01-00000F00","firmCode":"NC0058900000","balance":1,"isin":"RU000A0JX0J2X"}`),
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
			body, detail, err := h.CreateSecurityLimit(tt.req)
			assert.Equal(t, tt.wantBody, body)
			assert.Contains(t, detail, tt.wantDetail)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
