package client

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestClient_L(t *testing.T) {
	t.Run("nil client", func(t *testing.T) {
		l := (*Client)(nil).L()
		assert.Equal(t, zap.L(), l)
	})
	t.Run("non nil client", func(t *testing.T) {
		log, err := zap.NewProduction()
		require.NoError(t, err)
		cli := &Client{
			log: log,
		}
		assert.Equal(t, log, cli.L())
	})
}

func TestClient_P(t *testing.T) {
	t.Run("nil client", func(t *testing.T) {
		l := (*Client)(nil).P()
		assert.Nil(t, l)
	})
	t.Run("non nil client", func(t *testing.T) {
		cli := &Client{
			pool: &pgxpool.Pool{},
		}
		assert.Empty(t, cli.P())
	})
}
