package collections

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestNewSet(t *testing.T) {
	set := NewSet[int]()
	if assert.NotNil(t, set) {
		assert.Equal(t, make(map[int]struct{}), set.data)
	}
}

func TestSet_Len(t *testing.T) {
	t.Run("nil set", func(t *testing.T) {
		set := (*Set[int])(nil)
		assert.Equal(t, 0, set.Len())
	})
	t.Run("newly initialized set", func(t *testing.T) {
		set := NewSet[int]()
		assert.Equal(t, 0, set.Len())
		set.Add(1)
		assert.Equal(t, 1, set.Len())
	})
}

func TestSet_Add(t *testing.T) {
	t.Run("non-nil set async", func(t *testing.T) {
		t.Parallel()
		item := 1
		set := NewSet[int]()
		assert.False(t, set.Contain(item))
		var wg1, wg2 sync.WaitGroup

		wg1.Add(20)
		for i := 0; i < 20; i++ {
			go func(t *testing.T) {
				set.Add(item)
				assert.Equal(t, 1, set.Len())
				assert.True(t, set.Contain(item))
				wg1.Done()
			}(t)
		}
		wg1.Wait()
		item2 := 2
		assert.False(t, set.Contain(item2))
		wg2.Add(20)
		for i := 0; i < 20; i++ {
			go func(t *testing.T) {
				set.Add(item2)
				assert.Equal(t, 2, set.Len())
				assert.True(t, set.Contain(item2))
				wg2.Done()
			}(t)
		}
		wg2.Wait()
	})
	t.Run("nil set", func(t *testing.T) {
		set := (*Set[string])(nil)
		assert.Panics(t, func() {
			set.Add("some item")
		})
	})
}

func TestSet_AddMany(t *testing.T) {
	t.Run("nil set", func(t *testing.T) {
		set := (*Set[string])(nil)
		assert.Panics(t, func() {
			set.AddMany("some item", "other item")
		})
	})
	t.Run("non nil set", func(t *testing.T) {
		t.Parallel()
		item := 1
		set := NewSet[int]()
		assert.False(t, set.Contain(item))
		var wg1, wg2 sync.WaitGroup

		wg1.Add(20)
		for i := 0; i < 20; i++ {
			go func(t *testing.T) {
				set.AddMany(item)
				assert.Equal(t, 1, set.Len())
				assert.True(t, set.Contain(item))
				wg1.Done()
			}(t)
		}
		wg1.Wait()
		item2 := 2
		assert.False(t, set.Contain(item2))
		wg2.Add(20)
		for i := 0; i < 20; i++ {
			go func(t *testing.T) {
				set.AddMany(item2, item)
				assert.Equal(t, 2, set.Len())
				assert.True(t, set.Contain(item2))
				wg2.Done()
			}(t)
		}
		wg2.Wait()
	})
}

func TestSet_Contain(t *testing.T) {
	t.Run("nil set", func(t *testing.T) {
		set := (*Set[string])(nil)
		assert.False(t, set.Contain(""))
	})
}

func TestSet_Items(t *testing.T) {
	t.Run("nil set", func(t *testing.T) {
		set := (*Set[int])(nil)
		assert.Nil(t, set.Items())
	})
	t.Run("non nil set", func(t *testing.T) {
		set := NewSet[int]()
		items := set.Items()
		if assert.NotNil(t, items) {
			assert.Equal(t, 0, len(items))
		}
		set.AddMany(1, 2, 3)
		items = set.Items()
		if assert.NotNil(t, items) {
			assert.Equal(t, 3, len(items))
		}
	})
}

func TestDistinct(t *testing.T) {
	tt := []struct {
		name  string
		items []any
		want  assert.BoolAssertionFunc
	}{
		{"no items", []any{}, assert.False},
		{"positive #1", []any{1, 2, 3, 4}, assert.True},
		{"negative #1", []any{1, 2, 2, 3, 4}, assert.False},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.want(t, Distinct(tc.items...))
		})
	}
}
