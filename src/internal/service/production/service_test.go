package production

import (
	"github.com/stretchr/testify/assert"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/service/production/mocks"
	"go.uber.org/zap"
	"testing"
)

func TestService_ImplementsInterface(t *testing.T) {
	assert.Implements(t, new(controller.Service), new(Service))
}

func TestNew(t *testing.T) {
	t.Run("nil logger", func(t *testing.T) {
		s, err := New(nil, &mocks.MockStore{})
		assert.Nil(t, s)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrNilReference)
		}
	})
	t.Run("nil storage", func(t *testing.T) {
		s, err := New(zap.L(), nil)
		assert.Nil(t, s)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrNilReference)
		}
	})
	t.Run("positive", func(t *testing.T) {
		s, err := New(zap.L(), &mocks.MockStore{})
		assert.NoError(t, err)
		if assert.NotNil(t, s) {
			assert.Implements(t, new(controller.Service), s)
		}
	})
}
