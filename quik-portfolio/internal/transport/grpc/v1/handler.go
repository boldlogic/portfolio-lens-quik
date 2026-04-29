package v1

import (
	"context"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	quikv1 "github.com/boldlogic/portfolio-lens-quik/proto/gen/go/quik/v1"
	"github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

type Service interface {
	GetMoneyLimits(ctx context.Context, date time.Time) ([]quik.MoneyLimit, error)
	GetSecurityLimits(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error)
	GetCurrentQuotes(ctx context.Context) ([]models.CurrentQuote, error)
	GetCurrentQuotesForKeys(ctx context.Context, keys []string) ([]models.CurrentQuote, error)
}

type Handler struct {
	quikv1.UnimplementedLimitsServiceServer
	service Service
	logger  *zap.Logger
}

func NewHandler(svc Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: svc,
		logger:  logger,
	}
}

func ptr(s string) *string {
	return &s
}
