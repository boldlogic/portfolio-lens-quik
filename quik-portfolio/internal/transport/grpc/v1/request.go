package v1

import (
	"fmt"
	"time"

	quikv1 "github.com/boldlogic/portfolio-lens-quik/proto/gen/go/quik-portfolio/v1"
)

type limitsRequestParams struct {
	LoadDate          time.Time
	Limit             uint32
	Offset            uint64
	ClientCodes       []string
	IncludeTotalCount bool
}

const (
	defaultPageLimit uint32 = 100
	maxPageLimit     uint32 = 500
)

func parseLimitsRequestParams(req *quikv1.LimitsRequest) (limitsRequestParams, error) {
	date, err := protoDateToTime(req.GetLoadDate())
	if err != nil {
		return limitsRequestParams{}, fmt.Errorf("%w", err)
	}

	clients := req.GetClientCodes()

	limit := req.GetLimit()
	if limit == 0 {
		limit = defaultPageLimit
	}
	if limit > maxPageLimit {
		limit = maxPageLimit
	}

	offset := req.GetOffset()
	count := req.GetIncludeTotalCount()

	return limitsRequestParams{
		LoadDate:          date,
		ClientCodes:       clients,
		Limit:             limit,
		Offset:            offset,
		IncludeTotalCount: count,
	}, nil

}
