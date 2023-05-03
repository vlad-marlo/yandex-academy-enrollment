package middleware

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"strconv"
)

const (
	queryLimitParamName  = "limit"
	queryOffsetParamName = "offset"
	paginationCtxKey     = "pagination_context_key"
)

// PaginationOpts encapsulates limit and offset from developer into private fields.
//
// Opts can be accessed by getters.
type PaginationOpts struct {
	limit  int
	offset int
}

func NewPaginationOpts(limit, offset string) (opts *PaginationOpts) {
	opts = &PaginationOpts{
		limit:  1,
		offset: 0,
	}
	var err error
	if limit != "" {
		opts.limit, err = strconv.Atoi(limit)
		if err != nil {
			opts.limit = 1
		}
	}
	if offset != "" {
		opts.offset, err = strconv.Atoi(offset)
		if err != nil {
			opts.offset = 0
		}
	}
	return opts
}

// Limit is limit getter.
func (opts *PaginationOpts) Limit() int {
	if opts == nil {
		zap.L().Warn("unexpected got nil pagination opts")
		return 1
	}
	return opts.limit
}

// Offset is offset getter.
func (opts *PaginationOpts) Offset() int {
	if opts == nil {
		zap.L().Warn("unexpected got nil pagination opts")
		return 0
	}
	return opts.offset
}

// Paginator checks are pagination options provided and passes opts threw request context with options if provided and
// with defaults if not respectively.
func Paginator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		opts := NewPaginationOpts(c.QueryParam(queryLimitParamName), c.QueryParam(queryOffsetParamName))
		RequestWithPaginationOpts(c, opts)
		return next(c)
	}
}

// RequestWithPaginationOpts adds options to echo context.
func RequestWithPaginationOpts(c echo.Context, opts *PaginationOpts) {
	c.Set(paginationCtxKey, opts)
}

// GetPaginationOptsFromRequest return options from echo context.
func GetPaginationOptsFromRequest(c echo.Context) *PaginationOpts {
	return c.Get(paginationCtxKey).(*PaginationOpts)
}
