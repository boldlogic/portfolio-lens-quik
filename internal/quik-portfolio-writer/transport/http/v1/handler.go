package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/boldlogic/portfolio-lens-quik/internal/models"
	errmodel "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/handler"
	"go.uber.org/zap"
)

type Handler struct {
	commonHandler handler.Adapter
	service       Service
	logger        *zap.Logger
	apiKey        string
}

func NewHandler(commonHandler handler.Adapter, svc Service, logger *zap.Logger, apiKey string) *Handler {
	return &Handler{
		commonHandler: commonHandler,
		service:       svc,
		logger:        logger,
		apiKey:        apiKey,
	}
}

func (h *Handler) auth(next handler.HandlerFunc) handler.HandlerFunc {
	return func(r *http.Request) (any, string, error) {
		if h.apiKey == "" {
			return nil, "", fmt.Errorf("нет ключа")
		}
		if r.Header.Get("X-API-KEY") != h.apiKey {
			return nil, "", errmodel.ErrUnauthorized
		}
		return next(r)
	}

}

func (h *Handler) Adapt(fn handler.HandlerFunc) http.HandlerFunc {
	return h.commonHandler.Adapt(fn)
}

type Service interface {
	UpsertLimit(ctx context.Context, limit models.LimitInput) error
	UpsertLimits(ctx context.Context, limits []models.LimitLine) error
}
