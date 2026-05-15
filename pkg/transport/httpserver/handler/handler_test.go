package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boldlogic/packages/transport/httputils"
	"github.com/boldlogic/packages/utils/converters"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestHandler_Adapt_Success(t *testing.T) {
	tests := []struct {
		name            string
		method          string
		data            any
		wantStatus      int
		wantBody        string
		wantContentType bool
	}{
		{
			name:            "GET_возвращает_200_и_JSON",
			method:          http.MethodGet,
			data:            map[string]string{"result": "ok"},
			wantStatus:      http.StatusOK,
			wantBody:        `{"result":"ok"}`,
			wantContentType: true,
		},
		{
			name:            "POST_возвращает_201_и_JSON",
			method:          http.MethodPost,
			data:            map[string]string{"result": "created"},
			wantStatus:      http.StatusCreated,
			wantBody:        `{"result":"created"}`,
			wantContentType: true,
		},
		{
			name:            "nil_данные_возвращают_204_без_тела",
			method:          http.MethodGet,
			data:            nil,
			wantStatus:      http.StatusNoContent,
			wantBody:        "",
			wantContentType: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler()
			req := httptest.NewRequest(tt.method, "/", nil)
			rec := httptest.NewRecorder()

			h.Adapt(func(_ *http.Request) (any, string, error) {
				return tt.data, "", nil
			}).ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Equal(t, tt.wantBody, rec.Body.String())
			assertContentType(t, rec, tt.wantContentType)
		})
	}
}

func TestHandler_Adapt_ErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		detail     string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "unsupported_media_type",
			err:        httputils.ErrUnsupportedMediaType,
			detail:     "неподдерживаемый Content-Type",
			wantStatus: http.StatusUnsupportedMediaType,
			wantBody:   `{"title":"UNSUPPORTED_MEDIA_TYPE","status":415,"detail":"неподдерживаемый Content-Type"}`,
		},
		{
			name:       "request_entity_too_large",
			err:        httputils.ErrRequestEntityTooLarge,
			detail:     "тело запроса слишком большое",
			wantStatus: http.StatusRequestEntityTooLarge,
			wantBody:   `{"title":"REQUEST_ENTITY_TOO_LARGE","status":413,"detail":"тело запроса слишком большое"}`,
		},
		{
			name:       "validation_error_from_model",
			err:        models.ErrValidation,
			detail:     "некорректные входные данные",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"title":"VALIDATION_ERROR","status":400,"detail":"некорректные входные данные"}`,
		},
		{
			name:       "validation_error_from_body_reading",
			err:        httputils.ErrReadingBody,
			detail:     "ошибка чтения тела запроса",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"title":"VALIDATION_ERROR","status":400,"detail":"ошибка чтения тела запроса"}`,
		},
		{
			name:       "validation_error_from_wrong_json",
			err:        converters.ErrWrongJSON,
			detail:     "некорректный формат JSON",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"title":"VALIDATION_ERROR","status":400,"detail":"некорректный формат JSON"}`,
		},
		{
			name:       "business_validation_error",
			err:        models.ErrBusinessValidation,
			detail:     "не принят код или параметры",
			wantStatus: http.StatusUnprocessableEntity,
			wantBody:   `{"title":"BUSINESS_VALIDATION_ERROR","status":422,"detail":"не принят код или параметры"}`,
		},
		{
			name:       "not_found",
			err:        models.ErrNotFound,
			detail:     "запись не найдена",
			wantStatus: http.StatusNotFound,
			wantBody:   `{"title":"NOT_FOUND","status":404,"detail":"запись не найдена"}`,
		},
		{
			name:       "conflict",
			err:        models.ErrConflict,
			detail:     "запись уже существует",
			wantStatus: http.StatusConflict,
			wantBody:   `{"title":"CONFLICT","status":409,"detail":"запись уже существует"}`,
		},
		{
			name:       "internal_error",
			err:        errors.New("ошибка БД"),
			detail:     "техническая деталь",
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"title":"SERVER_ERROR","status":500,"detail":"что-то пошло не так"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			h.Adapt(func(_ *http.Request) (any, string, error) {
				return nil, tt.detail, tt.err
			}).ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Equal(t, tt.wantBody, rec.Body.String())
			assertContentType(t, rec, true)
		})
	}
}

func assertContentType(t *testing.T, rec *httptest.ResponseRecorder, want bool) {
	t.Helper()

	contentTypes := rec.Result().Header.Values("Content-Type")
	if want {
		require.NotEmpty(t, contentTypes)
		require.Contains(t, contentTypes[0], "application/json")
		require.Contains(t, contentTypes[0], "charset=UTF-8")
		return
	}
	require.Empty(t, contentTypes)
}
