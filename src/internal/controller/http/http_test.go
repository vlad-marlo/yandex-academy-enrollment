package http

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller/mocks"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"testing"
)

func TestRun(t *testing.T) {
	srv := testServer(t, nil)
	srv.configure()
	ctx := context.Background()
	require.NoError(t, srv.Start(ctx))
	assert.NoError(t, srv.Stop(ctx))
}

func TestNew(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		srv, err := New(zap.L(), &config{}, &config{}, &mocks.MockService{})
		assert.NoError(t, err)
		if assert.NotNil(t, srv) {
			assert.Equal(t, zap.L(), srv.log)
			assert.Equal(t, &config{}, srv.cfg)
			assert.Equal(t, &config{}, srv.rateCfg)
		}
	})
	t.Run("nil logger", func(t *testing.T) {
		srv, err := New(nil, &config{}, &config{}, &mocks.MockService{})
		assert.Nil(t, srv)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrNilReference)
		}
	})
	t.Run("nil config", func(t *testing.T) {
		srv, err := New(zap.L(), nil, &config{}, &mocks.MockService{})
		assert.Nil(t, srv)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrNilReference)
		}
	})
	t.Run("nil rate config", func(t *testing.T) {
		srv, err := New(zap.L(), &config{}, nil, &mocks.MockService{})
		assert.Nil(t, srv)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrNilReference)
		}
	})
}

func TestConfig(t *testing.T) {
	c := &config{}
	assert.Equal(t, 10, c.Burst())
	assert.Equal(t, rate.Limit(10), c.Limit())
	assert.Equal(t, bindAddr, c.BindAddr())
}
