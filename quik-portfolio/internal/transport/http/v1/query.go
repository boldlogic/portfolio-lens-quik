package v1

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/boldlogic/packages/transport/httpserver/request"
	"github.com/boldlogic/packages/utils/dates"
)

const (
	clientCodesQuery    = "clientCodes"
	totalCountQueryFlag = "includeTotalCount"
	dateQuery           = "loadDate"
	currencyQuery       = "currency"
)

func extractClientsQueryParam(r *http.Request) []string {
	raw := r.URL.Query().Get(clientCodesQuery)
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}

	}
	if len(out) == 0 {
		return nil
	}

	return out
}

func extractDateQueryParam(r *http.Request) (time.Time, error) {
	dateReq := r.URL.Query().Get(dateQuery)
	return dates.ParseWithDefaultNow(dateReq, dates.ISODateFormat)
}

func extractTotalCountFlag(r *http.Request) (bool, error) {
	flag, err := request.ParseBoolQuery(r, totalCountQueryFlag, false)
	if err != nil {
		return false, fmt.Errorf("%w %s", err, totalCountQueryFlag)
	}

	return flag, nil

}

type limitsListQuery struct {
	Date              time.Time
	Limit             uint32
	Offset            uint64
	ClientCodes       []string
	IncludeTotalCount bool
}

func parseLimitsListQuery(r *http.Request) (limitsListQuery, error) {
	date, err := extractDateQueryParam(r)
	if err != nil {
		return limitsListQuery{}, fmt.Errorf("%w. Ожидается YYYY-MM-DD", err)
	}

	limit, offset, err := request.ParseListPagination(r)
	if err != nil {
		return limitsListQuery{}, err
	}

	clients := extractClientsQueryParam(r)

	countFlag, err := extractTotalCountFlag(r)
	if err != nil {
		return limitsListQuery{}, err
	}

	return limitsListQuery{
		Date:              date,
		Limit:             limit,
		Offset:            offset,
		ClientCodes:       clients,
		IncludeTotalCount: countFlag,
	}, nil
}
