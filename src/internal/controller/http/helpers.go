package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/datetime"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/fielderr"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const (
	queryLimitParamName  = "limit"
	queryOffsetParamName = "offset"
)

// respond writes data to response writer.
//
// Passing nil data will write text status code to w.
func (srv *Controller) checkErr(c echo.Context, msg string, err error, fields ...zap.Field) error {
	var fieldErr *fielderr.Error
	if errors.As(err, &fieldErr) {
		srv.log.Warn(msg, append(fieldErr.Fields(), fields...)...)
		return c.JSON(fieldErr.CodeHTTP(), fieldErr.Data())
	}
	srv.log.Warn(msg, append(fields, zap.Error(err))...)
	return c.JSON(http.StatusBadRequest, model.BadRequestResponse{})
}

func (srv *Controller) dateFromContext(c echo.Context, queryParamName string) (*datetime.Date, error) {
	s := c.QueryParam(queryParamName)
	if s == "" {
		return datetime.Today(), nil
	}
	return datetime.ParseDate(s)
}

// PaginationOpts encapsulates limit and offset from developer into private fields.
//
// Opts can be accessed by getters.
type PaginationOpts struct {
	limit  int
	offset int
}

// NewPaginationOpts returns PaginationOpts with provided limit and offset.
//
// By default, values are parsing int from respective strings.
// Is string is not parsable then will be set default value - for limit it is 1 and for offset it is 0.
func NewPaginationOpts(limit, offset string) (opts *PaginationOpts) {
	opts = &PaginationOpts{
		limit:  1,
		offset: 0,
	}
	var err error

	opts.limit, err = strconv.Atoi(limit)
	if err != nil {
		zap.L().Info("got error while parsing pagination opts: limit", zap.Error(err))
		opts.limit = 1
	}

	opts.offset, err = strconv.Atoi(offset)
	if err != nil {
		zap.L().Info("got error while parsing pagination opts: offset", zap.Error(err))
		opts.offset = 0
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

// GetPaginationOptsFromRequest return options from echo context.
func GetPaginationOptsFromRequest(c echo.Context) *PaginationOpts {
	opts := NewPaginationOpts(c.QueryParam(queryLimitParamName), c.QueryParam(queryOffsetParamName))
	return opts
}
