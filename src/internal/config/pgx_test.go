package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func unsetEnv(t testing.TB, key, val string) func() {
	t.Helper()
	before := os.Getenv(key)
	require.NoError(t, os.Setenv(key, val))
	return func() {
		require.NoError(t, os.Setenv(key, before))
	}
}

func TestPgConfig(t *testing.T) {
	tt := []struct {
		name string
		cfg  *PgConfig
		want string
	}{
		{"nil cfg", nil, "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"},
		{"empty cfg", new(PgConfig), "postgresql://:@:0/?sslmode=disable"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.cfg.URI())
		})
	}
}

func TestNewPgConfig(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cfg, err := NewPgConfig()
		assert.NoError(t, err)
		assert.NotEmpty(t, cfg)
	})
	t.Run("negative", func(t *testing.T) {
		td := unsetEnv(t, "POSTGRES_PORT", "non int")
		cfg, err := NewPgConfig()
		assert.Error(t, err)
		assert.Nil(t, cfg)
		td()
		cfg, err = NewPgConfig()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
	})
}
