package v1

import (
	"context"
	"errors"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/models"
	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	quikv1 "github.com/boldlogic/portfolio-lens-quik/proto/gen/go/quik/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func moneyLimitToProto(l models.MoneyLimit) *quikv1.MoneyLimit {
	return &quikv1.MoneyLimit{
		ClientCode:   l.ClientCode,
		Currency:     l.Currency,
		Balance:      l.Balance.String(),
		FirmCode:     l.FirmCode,
		FirmName:     l.FirmName,
		PositionCode: l.PositionCode,
		SettleCode:   string(l.SettleCode),
		SourceDate:   timestamppb.New(l.SourceDate),
		LoadDate:     timestamppb.New(l.LoadDate),
	}
}

func (h *Handler) GetMoneyLimits(ctx context.Context, req *quikv1.GetMoneyLimitsRequest) (*quikv1.GetMoneyLimitsResponse, error) {
	date := time.Now()
	if req.GetDate() != nil {
		date = req.GetDate().AsTime()
	}

	limits, err := h.service.GetMoneyLimits(ctx, date)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		}
		h.logger.Error("лимиты ДС: чтение", zap.Error(err), zap.Time("date", date))
		return nil, status.Errorf(codes.Internal, "не удалось отправить money limits")
	}

	pbLimits := make([]*quikv1.MoneyLimit, len(limits))
	for i, l := range limits {
		pbLimits[i] = moneyLimitToProto(l)
	}

	return &quikv1.GetMoneyLimitsResponse{Limits: pbLimits}, nil
}

func (h *Handler) StreamMoneyLimits(req *quikv1.GetMoneyLimitsRequest, stream quikv1.LimitsService_StreamMoneyLimitsServer) error {
	date := time.Now()
	if req.GetDate() != nil {
		date = req.GetDate().AsTime()
	}

	limits, err := h.service.GetMoneyLimits(stream.Context(), date)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return status.Errorf(codes.InvalidArgument, "%v", err)
		}
		h.logger.Error("лимиты ДС: чтение", zap.Error(err), zap.Time("date", date))
		return status.Errorf(codes.Internal, "не удалось отправить money limits")
	}

	for _, l := range limits {
		if err := stream.Send(moneyLimitToProto(l)); err != nil {
			h.logger.Error("лимиты ДС: стрим, отправка", zap.Error(err))
			return status.Errorf(codes.Internal, "не удалось отправить money limit: %v", err)
		}
	}

	return nil
}
