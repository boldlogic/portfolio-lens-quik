package v1

import (
	"context"
	"errors"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	quikv1 "github.com/boldlogic/portfolio-lens-quik/proto/gen/go/quik-portfolio/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func moneyLimitToProto(l quik.MoneyLimit) *quikv1.MoneyLimit {
	return &quikv1.MoneyLimit{
		ClientCode:   l.ClientCode,
		Currency:     l.CurrencyCode,
		Balance:      l.Balance.String(),
		FirmCode:     l.FirmCode,
		FirmName:     l.FirmName,
		PositionCode: l.PositionCode,
		SettleCode:   string(l.SettleCode),
		SourceDate:   timeToProtoDate(l.SourceDate),
		LoadDate:     timeToProtoDate(l.LoadDate),
	}
}

func moneyLimitsToProto(l []quik.MoneyLimit) []*quikv1.MoneyLimit {
	out := make([]*quikv1.MoneyLimit, 0, len(l))
	for _, o := range l {
		out = append(out, moneyLimitToProto(o))
	}
	return out
}

func (h *Handler) GetMoneyLimits(ctx context.Context, req *quikv1.LimitsRequest) (*quikv1.GetMoneyLimitsResponse, error) {

	r, err := parseLimitsRequestParams(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)

	}

	limits, total, err := h.service.GetMoneyLimitsWithFilters(ctx, r.LoadDate, r.Limit, r.Offset, r.ClientCodes, r.IncludeTotalCount)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		}
		return nil, status.Errorf(codes.Internal, "что-то не так")
	}

	pbLimits := moneyLimitsToProto(limits)
	if r.IncludeTotalCount && total == nil {
		var z uint64 = 0
		total = new(z)
	}

	return &quikv1.GetMoneyLimitsResponse{
		Limits: pbLimits,
		Pagination: &quikv1.Pagination{
			Limit:  r.Limit,
			Offset: r.Offset,
			Total:  total,
		}}, nil
}
