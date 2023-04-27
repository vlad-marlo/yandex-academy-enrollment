package middleware

import (
	"github.com/labstack/echo"
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

// Limit is limit getter.
func (opts *PaginationOpts) Limit() int {
	if opts == nil {
		zap.L().Warn("unexpected got nil pagination opts")
		return 0
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
		var err error
		opts := new(PaginationOpts)
		opts.limit = 1 // setting up default limit
		opts.offset, err = strconv.Atoi(c.QueryParam(queryOffsetParamName))
		if err != nil {
			RequestWithPaginationOpts(c, opts)
			return next(c)
		}
		opts.limit, err = strconv.Atoi(c.QueryParam(queryLimitParamName))
		if err != nil {
			return next(c)
		}
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
