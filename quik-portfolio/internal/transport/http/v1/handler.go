package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	"github.com/boldlogic/portfolio-lens-quik/pkg/transport/httpserver/handler"
	"github.com/shopspring/decimal"
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
	GetMoneyLimitsWithFilters(ctx context.Context, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.MoneyLimit, totalCount *uint64, err error)
	GetSecurityLimitsWithFilters(ctx context.Context, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.SecurityLimit, totalCount *uint64, err error)
	GetSecurityLimitsOtcWithFilters(ctx context.Context, date time.Time, limit uint32, offset uint64, clientCodes []string, includeTotalCount bool) (result []quik.SecurityLimit, totalCount *uint64, err error)
	// GetPortfolio(ctx context.Context, targetCcy string) ([]quik.PortfolioEntry, error)
	GetMoneyPositions(ctx context.Context, currency *string) ([]quik.Position, decimal.Decimal, string, error)
}
