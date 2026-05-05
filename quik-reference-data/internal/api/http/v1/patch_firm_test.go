package v1

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/handler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const exampleURL = "/api/v1/quik/firms/"

type writeSvc struct {
	err error
}

func (s writeSvc) CreateFirm(ctx context.Context, code string, name string) (quik.Firm, error) {
	return quik.Firm{
		Id:   1,
		Code: code,
		Name: name,
	}, s.err
}

func (s writeSvc) UpdateFirm(ctx context.Context, id uint8, name string) (quik.Firm, error) {
	return quik.Firm{
		Id:   id,
		Code: "code",
		Name: name,
	}, s.err
}

func TestUpdateFirm(t *testing.T) {
	h := NewHandler(handler.NewHandler(), readSvc{}, writeSvc{}, zap.NewNop())
	t.Run("успешный_запрос", func(t *testing.T) {

		req := reqFirmWithID(http.MethodPatch, "1", bytes.NewBufferString(`{"firmName":"Новый брокер"}`))
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		body, detail, err := h.UpdateFirm(req)
		assert.Equal(t, "", detail)
		assert.NoError(t, err)
		assert.Equal(t, firmRespDTO{Id: 1, Code: "code", Name: "Новый брокер"}, body)
	})
	t.Run("некорректный_строковый_id", func(t *testing.T) {
		req := reqFirmWithID(http.MethodPatch, "a", bytes.NewBufferString(`{"firmName":"Новый брокер"}`))
		body, detail, err := h.UpdateFirm(req)
		assert.Nil(t, body)
		assert.Contains(t, detail, "некорректный id фирмы")
		assert.ErrorIs(t, err, models.ErrValidation)

	})
	t.Run("id>uint8", func(t *testing.T) {
		req := reqFirmWithID(http.MethodPatch, "1000", bytes.NewBufferString(`{"firmName":"Новый брокер"}`))
		body, detail, err := h.UpdateFirm(req)
		assert.Nil(t, body)
		assert.Contains(t, detail, "некорректный id фирмы")
		assert.ErrorIs(t, err, models.ErrValidation)

	})
	t.Run("id<0", func(t *testing.T) {
		req := reqFirmWithID(http.MethodPatch, "-1", bytes.NewBufferString(`{"firmName":"Новый брокер"}`))
		body, detail, err := h.UpdateFirm(req)
		assert.Nil(t, body)
		assert.Contains(t, detail, "некорректный id фирмы")
		assert.ErrorIs(t, err, models.ErrValidation)

	})
	t.Run("id=255", func(t *testing.T) {
		req := reqFirmWithID(http.MethodPatch, "255", bytes.NewBufferString(`{"firmName":"Новый брокер"}`))
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		body, detail, err := h.UpdateFirm(req)
		assert.Equal(t, "", detail)
		assert.NoError(t, err)
		assert.Equal(t, firmRespDTO{Id: 255, Code: "code", Name: "Новый брокер"}, body)

	})

}
