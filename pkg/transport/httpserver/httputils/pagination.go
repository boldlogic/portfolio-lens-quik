package httputils

import (
	"errors"
	"net/http"
	"strconv"
)

const (
	DefaultLimit = 50
	MaxLimit     = 200
)

var (
	ErrInvalidLimit  = errors.New("некорректный limit")
	ErrInvalidOffset = errors.New("некорректный offset")
)

func ParseListPagination(r *http.Request) (limit, offset int, err error) {
	limit = DefaultLimit
	offset = 0
	if v := r.URL.Query().Get("limit"); v != "" {
		n, e := strconv.Atoi(v)
		if e != nil || n < 1 {
			return 0, 0, ErrInvalidLimit
		}
		limit = n
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		n, e := strconv.Atoi(v)
		if e != nil || n < 0 {
			return 0, 0, ErrInvalidOffset
		}
		offset = n
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return limit, offset, nil
}
