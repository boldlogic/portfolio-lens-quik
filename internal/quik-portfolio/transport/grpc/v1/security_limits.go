package v1

import (
	"context"
	"errors"
	"strings"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	"github.com/boldlogic/portfolio-lens-quik/pkg/quik"
	quikv1 "github.com/boldlogic/portfolio-lens-quik/proto/gen/go/quik-portfolio/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func securityLimitToProto(l quik.SecurityLimit) *quikv1.SecurityLimit {
	pb := &quikv1.SecurityLimit{
		ClientCode:              l.ClientCode,
		SecCode:                 l.SecCode,
		Balance:                 l.Balance.String(),
		AcquisitionCurrencyCode: l.AcquisitionCurrencyCode,
		FirmCode:       l.FirmCode,
		FirmName:       l.FirmName,
		TradeAccount:   l.TradeAccount,
		SettleCode:     string(l.SettleCode),
		SourceDate:     timeToProtoDate(l.SourceDate),
		LoadDate:       timeToProtoDate(l.LoadDate),
	}
	if l.ISIN != "" {
		pb.Isin = &l.ISIN
	}
	if sn := strings.TrimSpace(l.ShortName); sn != "" {
		pb.ShortName = &sn
	}
	return pb
}

func securityLimitsToProto(l []quik.SecurityLimit) []*quikv1.SecurityLimit {
	out := make([]*quikv1.SecurityLimit, 0, len(l))
	for _, o := range l {
		out = append(out, securityLimitToProto(o))
	}
	return out
}
func (h *Handler) GetSecurityLimits(ctx context.Context, req *quikv1.LimitsRequest) (*quikv1.GetSecurityLimitsResponse, error) {
	r, err := parseLimitsRequestParams(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)

	}

	limits, total, err := h.service.GetSecurityLimitsWithFilters(ctx, r.LoadDate, r.Limit, r.Offset, r.ClientCodes, r.IncludeTotalCount)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		}
		return nil, status.Errorf(codes.Internal, "не удалось получить security limits")
	}
	if r.IncludeTotalCount && total == nil {
		var z uint64 = 0
		total = new(z)
	}

	pbLimits := securityLimitsToProto(limits)

	return &quikv1.GetSecurityLimitsResponse{
		Limits: pbLimits,
		Pagination: &quikv1.Pagination{
			Limit:  r.Limit,
			Offset: r.Offset,
			Total:  total,
		}}, nil
}
