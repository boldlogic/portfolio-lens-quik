package v1

import (
	"context"
	"net/http"
	"time"

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
	GetMoneyLimits(ctx context.Context, date time.Time) ([]quik.MoneyLimit, error)
	GetSecurityLimits(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error)
	GetSecurityLimitsOtc(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error)

	CreateMoneyLimit(ctx context.Context, ml quik.MoneyLimit) (quik.MoneyLimit, error)
	CreateSecurityLimit(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error)
	CreateSecurityLimitOtc(ctx context.Context, sec quik.SecurityLimit) (quik.SecurityLimit, error)

	GetLimits(ctx context.Context, date time.Time) ([]quik.Limit, error)
	GetPortfolio(ctx context.Context, targetCcy string) ([]quik.PortfolioEntry, error)
}
