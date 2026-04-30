package v1

import (
	"context"
	"net/http"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/handler"
	"go.uber.org/zap"
)

type Handler struct {
	commonHandler handler.Adapter
	readService   ReadService
	writeService  WriteService
	logger        *zap.Logger
}

type ReadService interface {
	GetFirms(ctx context.Context) ([]quik.Firm, error)
	GetFirmByID(ctx context.Context, id uint8) (quik.Firm, error)
}

type WriteService interface {
	CreateFirm(ctx context.Context, code string, name string) (quik.Firm, error)
	UpdateFirm(ctx context.Context, id uint8, name string) (quik.Firm, error)
}

func NewHandler(
	commonHandler handler.Adapter,
	readService ReadService,
	writeService WriteService,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		commonHandler: commonHandler,
		readService:   readService,
		writeService:  writeService,
		logger:        logger,
	}
}

func (h *Handler) Adapt(fn handler.HandlerFunc) http.HandlerFunc {
	return h.commonHandler.Adapt(fn)
}
