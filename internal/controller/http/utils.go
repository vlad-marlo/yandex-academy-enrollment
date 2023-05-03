package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/fielderr"
)

// respond writes data to response writer.
//
// Passing nil data will write text status code to w.
func (srv *Controller) checkErr(c echo.Context, err error) error {
	var fieldErr *fielderr.Error
	if errors.As(err, &fieldErr) {
		if fieldErr.Data() != nil {
			return c.JSON(fieldErr.CodeHTTP(), fieldErr.Data())
		}
		return c.NoContent(fieldErr.CodeHTTP())
	}
	return err
}
