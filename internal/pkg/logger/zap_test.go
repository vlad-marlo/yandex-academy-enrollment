package logger

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestNew_Positive(t *testing.T) {
	global := zap.L()
	l, err := New()
	assert.NoError(t, err)
	// check that logger replacing global logger
	assert.Equal(t, l, zap.L())
	assert.NotEqual(t, global, zap.L())
}

func TestNew_Negative(t *testing.T) {

}
