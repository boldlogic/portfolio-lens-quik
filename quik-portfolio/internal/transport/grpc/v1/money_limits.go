package v1

import (
	"context"
	"errors"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
	quikv1 "github.com/boldlogic/portfolio-lens-quik/proto/gen/go/quik-portfolio/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func moneyLimitToProto(l quik.MoneyLimit) *quikv1.MoneyLimit {
	return &quikv1.MoneyLimit{
		ClientCode:   l.ClientCode,
		Currency:     l.Currency,
		Balance:      l.Balance.String(),
		FirmCode:     l.FirmCode,
		FirmName:     l.FirmName,
		PositionCode: l.PositionCode,
		SettleCode:   string(l.SettleCode),
		SourceDate:   timeToProtoDate(l.SourceDate),
		LoadDate:     timeToProtoDate(l.LoadDate),
	}
}

func moneyLimitsToResp(l []quik.MoneyLimit) []*quikv1.MoneyLimit {
	out := make([]*quikv1.MoneyLimit, 0, len(l))
	for _, o := range l {
		out = append(out, moneyLimitToProto(o))
	}
	return out
}

func (h *Handler) GetMoneyLimits(ctx context.Context, req *quikv1.GetMoneyLimitsRequest) (*quikv1.GetMoneyLimitsResponse, error) {
	date, err := protoDateToTime(req.GetDate())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	limits, err := h.service.GetMoneyLimits(ctx, date)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		}
		h.logger.Error("лимиты ДС: чтение gRPC", zap.Error(err), zap.Time("date", date))
		return nil, status.Errorf(codes.Internal, "не удалось отправить money limits")
	}

	pbLimits := moneyLimitsToResp(limits)

	return &quikv1.GetMoneyLimitsResponse{Limits: pbLimits}, nil
}
