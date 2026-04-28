package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/models"
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
	GetMoneyLimits(ctx context.Context, date time.Time) ([]models.MoneyLimit, error)
	GetSecurityLimits(ctx context.Context, date time.Time) ([]models.SecurityLimit, error)
	GetSecurityLimitsOtc(ctx context.Context, date time.Time) ([]models.SecurityLimit, error)

	CreateMoneyLimit(ctx context.Context, ml models.MoneyLimit) (models.MoneyLimit, error)
	CreateSecurityLimit(ctx context.Context, sec models.SecurityLimit) (models.SecurityLimit, error)
	CreateSecurityLimitOtc(ctx context.Context, sec models.SecurityLimit) (models.SecurityLimit, error)

	GetLimits(ctx context.Context, date time.Time) ([]models.Limit, error)
	GetPortfolio(ctx context.Context, targetCcy string) ([]models.PortfolioEntry, error)
	GetFirms(ctx context.Context) ([]quik.Firm, error)
	GetFirmByID(ctx context.Context, id uint8) (quik.Firm, error)
	CreateFirm(ctx context.Context, code string, name string) (quik.Firm, error)
	UpdateFirm(ctx context.Context, id uint8, name string) (quik.Firm, error)
}
