package collections

import (
	"sync"
)

// Set is implementation of std container.
//
// Set stores distinct items always only in one
type Set[T comparable] struct {
	mu   sync.RWMutex
	data map[T]struct{}
}

// NewSet returns set implementation.
func NewSet[T comparable](items ...T) *Set[T] {
	s := &Set[T]{
		mu:   sync.RWMutex{},
		data: make(map[T]struct{}),
	}
	s.AddMany(items...)
	return s
}

// Add adds item to Set.
//
// Adding is safe to use in concurrent code.
func (set *Set[T]) Add(item T) {
	if set == nil {
		panic("set must be non-nil object")
	}
	set.mu.Lock()
	set.data[item] = struct{}{}
	set.mu.Unlock()
}

// AddMany stores all items from slice in Set.
func (set *Set[T]) AddMany(data ...T) {
	if set == nil {
		panic("set must be non-nil object")
	}
	set.mu.Lock()
	for _, item := range data {
		set.data[item] = struct{}{}
	}
	set.mu.Unlock()
}

// Len return count of distinct items, stored in Set.
func (set *Set[T]) Len() int {
	if set == nil {
		return 0
	}
	set.mu.RLock()
	res := len(set.data)
	set.mu.RUnlock()
	return res
}

// Items returns slice of items stored in Set.
func (set *Set[T]) Items() []T {
	if set == nil {
		return nil
	}
	set.mu.RLock()
	res := make([]T, 0, len(set.data))
	for k := range set.data {
		res = append(res, k)
	}
	set.mu.RUnlock()
	return res
}

// Contain returns is item stored in set.
//
// Return true if item is stored and false if not respectively.
func (set *Set[T]) Contain(item T) bool {
	if set == nil {
		return false
	}
	set.mu.RLock()
	_, ok := set.data[item]
	set.mu.RUnlock()
	return ok
}

// Distinct return are items distinct or not.
func Distinct[T comparable](items ...T) bool {
	return len(items) == NewSet[T](items...).Len() && len(items) != 0
}
