package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/datetime"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/fielderr"
	"go.uber.org/zap"
)

// respond writes data to response writer.
//
// Passing nil data will write text status code to w.
func (srv *Controller) checkErr(c echo.Context, err error) error {
	var fieldErr *fielderr.Error
	if errors.As(err, &fieldErr) {
		return c.JSON(fieldErr.CodeHTTP(), fieldErr.Data())
	}
	zap.L().Warn("checked error", zap.Error(err))
	return err
}

func (srv *Controller) dateFromContext(c echo.Context, queryParamName string) (*datetime.Date, error) {
	s := c.QueryParam(queryParamName)
	if s == "" {
		return datetime.Today(), nil
	}
	return datetime.ParseDate(s)
}
