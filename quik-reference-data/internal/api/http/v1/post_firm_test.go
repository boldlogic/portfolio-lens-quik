package v1

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boldlogic/packages/transport/httputils"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/handler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateFirm(t *testing.T) {
	h := NewHandler(handler.NewHandler(), readSvc{}, writeSvc{}, zap.NewNop())
	t.Run("успешный_запрос", func(t *testing.T) {

		payload := bytes.NewBufferString(`{"firmName":"Новый брокер","firmCode":"code"}`)
		req := httptest.NewRequest(http.MethodPost, "/", payload)
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		body, detail, err := h.CreateFirm(req)
		assert.Equal(t, "", detail)
		assert.NoError(t, err)
		assert.Equal(t, firmRespDTO{Id: 1, Code: "code", Name: "Новый брокер"}, body)
	})
	t.Run("пропущен_firmName", func(t *testing.T) {
		payload := bytes.NewBufferString(`{"firmCode":"code"}`)
		req := httptest.NewRequest(http.MethodPost, "/", payload)
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		body, detail, err := h.CreateFirm(req)
		assert.Nil(t, body)
		assert.NotEmpty(t, detail)
		assert.ErrorIs(t, err, models.ErrValidation)
	})
	t.Run("пропущен_firmCode", func(t *testing.T) {
		payload := bytes.NewBufferString(`{"firmName":"Новый брокер"`)
		req := httptest.NewRequest(http.MethodPost, "/", payload)
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		body, detail, err := h.CreateFirm(req)
		assert.Nil(t, body)
		assert.NotEmpty(t, detail)
		assert.ErrorIs(t, err, models.ErrValidation)
	})
	t.Run("UnsupportedMediaType", func(t *testing.T) {
		payload := bytes.NewBufferString(`{"firmName":"Новый брокер","firmCode":"code"}`)
		req := httptest.NewRequest(http.MethodPost, "/", payload)
		body, detail, err := h.CreateFirm(req)
		assert.Nil(t, body)
		assert.NotEmpty(t, detail)
		assert.ErrorIs(t, err, httputils.ErrUnsupportedMediaType)
	})

	t.Run("firmCode>ограничения", func(t *testing.T) {
		payload := bytes.NewBufferString(`{"firmName":"Новый брокер","firmCode":"code567890123"}`)
		req := httptest.NewRequest(http.MethodPost, "/", payload)
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		body, detail, err := h.CreateFirm(req)
		assert.Nil(t, body)
		assert.NotEmpty(t, detail)
		assert.ErrorIs(t, err, models.ErrValidation)
	})
}
