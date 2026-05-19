package v1

import (
	"context"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	quikv1 "github.com/boldlogic/portfolio-lens-quik/proto/gen/go/quik-portfolio/v1"
	"go.uber.org/zap"
)

type Service interface {
	GetMoneyLimits(ctx context.Context, date time.Time) ([]quik.MoneyLimit, error)
	GetSecurityLimits(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error)
	GetSecurityLimitsOtc(ctx context.Context, date time.Time) ([]quik.SecurityLimit, error)
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

// func ptr(s string) *string {
// 	return &s
// }
