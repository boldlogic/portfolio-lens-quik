package v1

import (
	"fmt"
	"time"

	quikv1 "github.com/boldlogic/portfolio-lens-quik/proto/gen/go/quik-portfolio/v1"
)

type limitsListQuery struct {
	Date              time.Time
	Limit             uint32
	Offset            uint64
	ClientCodes       []string
	IncludeTotalCount bool
}

const (
	defaultLimit uint32 = 100
	maxLimit     uint32 = 500
)

func extractReqFields(req *quikv1.LimitsRequest) (limitsListQuery, error) {
	date, err := protoDateToTime(req.GetLoadDate())
	if err != nil {
		return limitsListQuery{}, fmt.Errorf("%w", err)
	}

	clients := req.GetClientCodes()

	limit := req.GetLimit()
	if limit > maxLimit {
		return limitsListQuery{}, fmt.Errorf("некорректный limit %d", limit)
	}
	if limit == 0 {
		limit = defaultLimit
	}

	offset := req.GetOffset()
	// if offset < 0 {
	// 	return limitsListQuery{}, fmt.Errorf("некорректный offset %d", offset)
	// }
	count := req.GetIncludeTotalCount()

	return limitsListQuery{
		Date:              date,
		ClientCodes:       clients,
		Limit:             limit,
		Offset:            offset,
		IncludeTotalCount: count,
	}, nil

}
