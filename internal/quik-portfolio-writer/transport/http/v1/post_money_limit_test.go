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

func TestCreateMoneyLimit(t *testing.T) {
	t.Parallel()

	loadDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)
	sourceDate := time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local)
	internalSvcErr := errors.New("временная_ошибка_хранилища")

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
				"clientCode":"AB12CD",
				"currency":"RUB",
				"positionCode":"EQTV",
				"settleCode":"Tx",
				"firmCode": "NC0058900000",
				"balance":"331.10"
			}`),
			svc: svc{
				createMoneyLimit: func(ctx context.Context, ml quik.MoneyLimit) (quik.MoneyLimit, error) {
					ml.LoadDate = loadDate
					ml.SourceDate = sourceDate
					ml.FirmCode = "NC0058900000"
					ml.FirmName = "Фирма брокера"
					return ml, nil
				},
			},
			wantBody: moneyLimitDTO{
				LoadDate:     "2025-01-01",
				SourceDate:   "2025-01-02",
				ClientCode:   "AB12CD",
				Currency:     "RUB",
				PositionCode: "EQTV",
				SettleCode:   "Tx",
				FirmCode:     "NC0058900000",
				FirmName:     "Фирма брокера",
				Balance:      decimal.RequireFromString("331.10"),
			},
		},
		{
			name:       "пустой_clientCode",
			req:        reqJSON(`{"clientCode":"","currency":"RUB","firmCode": "NC0058900000"}`),
			wantBody:   nil,
			wantDetail: "clientCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "clientCode>12",
			req:        reqJSON(`{"clientCode":"1234567890123","currency":"RUB","firmCode": "NC0058900000"}`),
			wantBody:   nil,
			wantDetail: "clientCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "пустой_currency",
			req:        reqJSON(`{"clientCode":"TBANK","currency":"","firmCode": "NC0058900000"}`),
			wantBody:   nil,
			wantDetail: "currency",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "пустой_firmCode",
			req:        reqJSON(`{"clientCode":"TBANK","currency":"RUR","firmCode": ""}`),
			wantBody:   nil,
			wantDetail: "firmCode",
			wantErr:    models.ErrValidation,
		},
		{
			name: "UnsupportedMediaType",
			req: httptest.NewRequest(http.MethodPost, exampleURL, bytes.NewBufferString(
				`{"clientCode":"AB12CD","currency":"RUB","firmCode":"NC0058900000","balance":"1"}`,
			)),
			wantBody:   nil,
			wantDetail: "Content-Type",
			wantErr:    httputils.ErrUnsupportedMediaType,
		},
		{
			name:       "битый_json",
			req:        reqJSON(`{"clientCode":"TBANK","currency":"RUR","firmCode":"NC0058900000","balance":"1"`),
			wantBody:   nil,
			wantDetail: "",
			wantErr:    models.ErrValidation,
		},
		{
			name: "конфликт_ключа",
			req: reqJSON(`{
				"clientCode":"AB12CD",
				"currency":"RUB",
				"firmCode":"NC0058900000",
				"balance":"331.10"
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
				"currency":"RUB",
				"firmCode":"NC0058900000",
				"balance":"331.10"
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
				"currency":"RUB",
				"firmCode":"NC0058900000",
				"balance":"1"
			}`),
			svc:        svc{err: internalSvcErr},
			wantBody:   nil,
			wantDetail: "",
			wantErr:    internalSvcErr,
		},
		{
			name:       "currency_длина_>3",
			req:        reqJSON(`{"clientCode":"TBANK","currency":"RUBB","firmCode":"NC0058900000","balance":"1"}`),
			wantBody:   nil,
			wantDetail: "currency",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "firmCode_длина_>12",
			req:        reqJSON(`{"clientCode":"TBANK","currency":"RUB","firmCode":"NC00589000001","balance":"1"}`),
			wantBody:   nil,
			wantDetail: "firmCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "positionCode_длина_>4",
			req:        reqJSON(`{"clientCode":"TBANK","currency":"RUB","positionCode":"EQTVX","firmCode":"NC0058900000","balance":"1"}`),
			wantBody:   nil,
			wantDetail: "positionCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "settleCode_длина_>5",
			req:        reqJSON(`{"clientCode":"TBANK","currency":"RUB","settleCode":"123456","firmCode":"NC0058900000","balance":"1"}`),
			wantBody:   nil,
			wantDetail: "settleCode",
			wantErr:    models.ErrValidation,
		},
		{
			name:       "некорректный_balance",
			req:        reqJSON(`{"clientCode":"TBANK","currency":"RUB","firmCode":"NC0058900000","balance":"not_a_number"}`),
			wantBody:   nil,
			wantDetail: "некорректный формат JSON",
			wantErr:    models.ErrValidation,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := newTestHandler(tt.svc)
			body, detail, err := h.postMoneyLimit(tt.req)
			assert.Equal(t, tt.wantBody, body)
			assert.Contains(t, detail, tt.wantDetail)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
