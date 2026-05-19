package v1

import (
	"context"
	"errors"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	quikv1 "github.com/boldlogic/portfolio-lens-quik/proto/gen/go/quik-portfolio/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetSecurityLimitsOtc(ctx context.Context, req *quikv1.GetSecurityLimitsRequest) (*quikv1.GetSecurityLimitsResponse, error) {
	date, err := protoDateToTime(req.GetDate())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	limits, err := h.service.GetSecurityLimitsOtc(ctx, date)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		}
		h.logger.Error("лимиты бумаги OTC: чтение gRPC", zap.Error(err), zap.Time("date", date))
		return nil, status.Errorf(codes.Internal, "не удалось получить security limits OTC")
	}

	return &quikv1.GetSecurityLimitsResponse{Limits: securityLimitsToResp(limits)}, nil
}
