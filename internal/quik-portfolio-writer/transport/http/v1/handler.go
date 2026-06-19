package v1

import (
	"context"
	"net/http"

	"github.com/boldlogic/portfolio-lens-quik/internal/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/handler"
	"go.uber.org/zap"
)

type Handler struct {
	commonHandler handler.Adapter
	service       Service
	logger        *zap.Logger
}

func NewHandler(commonHandler handler.Adapter, svc Service, logger *zap.Logger) *Handler {
	return &Handler{
		commonHandler: commonHandler,
		service:       svc,
		logger:        logger,
	}
}

func (h *Handler) Adapt(fn handler.HandlerFunc) http.HandlerFunc {
	return h.commonHandler.Adapt(fn)
}

type Service interface {
	CreateLimit(ctx context.Context, l quik.Limit) (quik.Limit, error)
	UpsertLimits(ctx context.Context, limits []models.LimitLine) error
}
