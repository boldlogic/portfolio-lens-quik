package v1

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/models"
	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	quikv1 "github.com/boldlogic/portfolio-lens-quik/proto/gen/go/quik/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func securityLimitToProto(l models.SecurityLimit) *quikv1.SecurityLimit {
	pb := &quikv1.SecurityLimit{
		ClientCode:     l.ClientCode,
		Ticker:         l.Ticker,
		Balance:        l.Balance.String(),
		AcquisitionCcy: l.AcquisitionCcy,
		FirmCode:       l.FirmCode,
		FirmName:       l.FirmName,
		TradeAccount:   l.TradeAccount,
		SettleCode:     string(l.SettleCode),
		SourceDate:     timestamppb.New(l.SourceDate),
		LoadDate:       timestamppb.New(l.LoadDate),
	}
	if l.ISIN != nil {
		pb.Isin = l.ISIN
	}
	if sn := strings.TrimSpace(l.ShortName); sn != "" {
		pb.ShortName = ptr(sn)
	}
	return pb
}

func (h *Handler) GetSecurityLimits(ctx context.Context, req *quikv1.GetSecurityLimitsRequest) (*quikv1.GetSecurityLimitsResponse, error) {
	date := time.Now()
	if req.GetDate() != nil {
		date = req.GetDate().AsTime()
	}

	limits, err := h.service.GetSecurityLimits(ctx, date)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		}
		h.logger.Error("лимиты (бумаги): чтение", zap.Error(err), zap.Time("date", date))
		return nil, status.Errorf(codes.Internal, "не удалось получить security limits")
	}

	pbLimits := make([]*quikv1.SecurityLimit, len(limits))
	for i, l := range limits {
		pbLimits[i] = securityLimitToProto(l)
	}

	return &quikv1.GetSecurityLimitsResponse{Limits: pbLimits}, nil
}

func (h *Handler) StreamSecurityLimits(req *quikv1.GetSecurityLimitsRequest, stream quikv1.LimitsService_StreamSecurityLimitsServer) error {
	date := time.Now()
	if req.GetDate() != nil {
		date = req.GetDate().AsTime()
	}

	limits, err := h.service.GetSecurityLimits(stream.Context(), date)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return status.Errorf(codes.InvalidArgument, "%v", err)
		}
		h.logger.Error("лимиты (бумаги): чтение", zap.Error(err), zap.Time("date", date))
		return status.Errorf(codes.Internal, "не удалось отправить security limits")
	}

	for _, l := range limits {
		if err := stream.Send(securityLimitToProto(l)); err != nil {
			h.logger.Error("лимиты (бумаги): стрим, отправка", zap.Error(err))
			return status.Errorf(codes.Internal, "failed to stream security limit: %v", err)
		}
	}

	return nil
}
