package main

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"testing"
)

func TestCreateApp(t *testing.T) {
	assert.NoError(t, fx.ValidateApp(CreateApp()))
}
