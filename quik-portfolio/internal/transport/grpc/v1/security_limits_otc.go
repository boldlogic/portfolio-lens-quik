package v1

import (
	"context"
	"errors"

	md "github.com/boldlogic/portfolio-lens-quik/pkg/models"
	quikv1 "github.com/boldlogic/portfolio-lens-quik/proto/gen/go/quik-portfolio/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetSecurityLimitsOtc(ctx context.Context, req *quikv1.LimitsRequest) (*quikv1.GetSecurityLimitsResponse, error) {
	r, err := parseLimitsRequestParams(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)

	}

	limits, total, err := h.service.GetSecurityLimitsOtcWithFilters(ctx, r.LoadDate, r.Limit, r.Offset, r.ClientCodes, r.IncludeTotalCount)
	if err != nil {
		if errors.Is(err, md.ErrBusinessValidation) {
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		}

		return nil, status.Errorf(codes.Internal, "не удалось получить security limits OTC")
	}

	if r.IncludeTotalCount && total == nil {
		var z uint64 = 0
		total = new(z)
	}

	return &quikv1.GetSecurityLimitsResponse{
		Limits: securityLimitsToProto(limits),
		Pagination: &quikv1.Pagination{
			Limit:  r.Limit,
			Offset: r.Offset,
			Total:  total,
		}}, nil
}
