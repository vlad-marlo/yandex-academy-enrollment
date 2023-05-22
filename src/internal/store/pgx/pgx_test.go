package pgx

import (
	"github.com/stretchr/testify/assert"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/pgx/client"
	"testing"
)

func TestNew_Positive(t *testing.T) {
	cli, td := client.NewTest(t)
	defer td()
	s, err := New(cli)
	if assert.NotNil(t, s) {
		assert.Equal(t, cli.L(), s.log)
		assert.Equal(t, cli.P(), s.pool)
	}
	assert.NoError(t, err)
}

func TestNew_Negative(t *testing.T) {
	s, err := New(nil)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrNilReference)
	}
	assert.Nil(t, s)
}
