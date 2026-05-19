package v1

import (
	"context"
	"errors"
	"strings"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	quikv1 "github.com/boldlogic/portfolio-lens-quik/proto/gen/go/quik-portfolio/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func securityLimitToProto(l quik.SecurityLimit) *quikv1.SecurityLimit {
	pb := &quikv1.SecurityLimit{
		ClientCode:     l.ClientCode,
		Ticker:         l.Ticker,
		Balance:        l.Balance.String(),
		AcquisitionCcy: l.AcquisitionCcy,
		FirmCode:       l.FirmCode,
		FirmName:       l.FirmName,
		TradeAccount:   l.TradeAccount,
		SettleCode:     string(l.SettleCode),
		SourceDate:     timeToProtoDate(l.SourceDate),
		LoadDate:       timeToProtoDate(l.LoadDate),
	}
	if l.ISIN != nil {
		pb.Isin = l.ISIN
	}
	if sn := strings.TrimSpace(l.ShortName); sn != "" {
		pb.ShortName = &sn
	}
	return pb
}

func securityLimitsToResp(l []quik.SecurityLimit) []*quikv1.SecurityLimit {
	out := make([]*quikv1.SecurityLimit, 0, len(l))
	for _, o := range l {
		out = append(out, securityLimitToProto(o))
	}
	return out
}
func (h *Handler) GetSecurityLimits(ctx context.Context, req *quikv1.GetSecurityLimitsRequest) (*quikv1.GetSecurityLimitsResponse, error) {
	date, err := protoDateToTime(req.GetDate())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	limits, err := h.service.GetSecurityLimits(ctx, date)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		}
		h.logger.Error("лимиты бумаги: чтение gRPC", zap.Error(err), zap.Time("date", date))
		return nil, status.Errorf(codes.Internal, "не удалось получить security limits")
	}

	pbLimits := securityLimitsToResp(limits)

	return &quikv1.GetSecurityLimitsResponse{Limits: pbLimits}, nil
}
