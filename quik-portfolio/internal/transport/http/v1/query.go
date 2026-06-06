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
	queryParamClientCodes       = "clientCodes"
	queryParamIncludeTotalCount = "includeTotalCount"
	queryParamLoadDate          = "loadDate"
	queryParamCurrency          = "currency"
)

func parseClientCodesQueryParam(r *http.Request) []string {
	raw := r.URL.Query().Get(queryParamClientCodes)
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

func parseLoadDateQueryParam(r *http.Request) (time.Time, error) {
	dateReq := r.URL.Query().Get(queryParamLoadDate)
	return dates.ParseWithDefaultNow(dateReq, dates.ISODateFormat)
}

func parseIncludeTotalCountQueryParam(r *http.Request) (bool, error) {
	flag, err := request.ParseBoolQuery(r, queryParamIncludeTotalCount, false)
	if err != nil {
		return false, fmt.Errorf("%w %s", err, queryParamIncludeTotalCount)
	}

	return flag, nil

}

func parseCurrencyQueryParam(r *http.Request) *string {
	targetCurrency := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get(queryParamCurrency)))
	if targetCurrency == "" {
		return nil
	}
	return &targetCurrency
}

type limitsQueryParams struct {
	LoadDate          time.Time
	Limit             uint32
	Offset            uint64
	ClientCodes       []string
	IncludeTotalCount bool
}

func parseLimitsQueryParams(r *http.Request) (limitsQueryParams, error) {
	date, err := parseLoadDateQueryParam(r)
	if err != nil {
		return limitsQueryParams{}, fmt.Errorf("%w. Ожидается YYYY-MM-DD", err)
	}

	limit, offset, err := request.ParseListPagination(r)
	if err != nil {
		return limitsQueryParams{}, err
	}

	clients := parseClientCodesQueryParam(r)

	countFlag, err := parseIncludeTotalCountQueryParam(r)
	if err != nil {
		return limitsQueryParams{}, err
	}

	return limitsQueryParams{
		LoadDate:          date,
		Limit:             limit,
		Offset:            offset,
		ClientCodes:       clients,
		IncludeTotalCount: countFlag,
	}, nil
}
