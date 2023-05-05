package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testConfig int

func (t testConfig) Limit() rate.Limit {
	return rate.Limit(t)
}

func (t testConfig) Burst() int {
	return int(t)
}

func TestRateLimiter(t *testing.T) {

	cfg := testConfig(1)
	mw := RateLimiter(cfg)
	h := mw(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	assert.NotNil(t, h)
	e := echo.New()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	defer assert.NoError(t, r.Body.Close())
	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())
	c := e.NewContext(r, w)
	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusOK, w.Code)
	}

	r2 := httptest.NewRequest(http.MethodGet, "/", nil)
	defer assert.NoError(t, r2.Body.Close())
	w2 := httptest.NewRecorder()
	defer assert.NoError(t, w2.Result().Body.Close())
	c2 := e.NewContext(r2, w2)
	if assert.NoError(t, h(c2)) {
		assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	}
}
