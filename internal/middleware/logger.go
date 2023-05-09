package middleware

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func LogRequest(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			dur := time.Now().Sub(start)
			log.Info(
				"got request",
				zap.Duration("duration", dur),
				zap.String("status text", http.StatusText(c.Response().Status)),
			)
			return err
		}
	}
}
