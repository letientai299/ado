// Package fp provides functional programming utilities.
package fp

import "slices"

func Map[I, O any](in []I, fn func(I) O) []O {
	if len(in) == 0 {
		return nil
	}
	result := make([]O, 0, len(in))
	for _, item := range in {
		result = append(result, fn(item))
	}
	return result
}

type Optional[T any] struct{ value *T }

func (o Optional[T]) IsSome() bool { return o.value != nil }
func (o Optional[T]) IsNil() bool  { return o.value == nil }
func (o Optional[T]) Get() T       { return *o.value }

func Nil[T any]() Optional[T]         { return Optional[T]{nil} }
func Some[T any](value T) Optional[T] { return Optional[T]{value: &value} }

func Not[E any](fn func(E) bool) func(E) bool {
	return func(value E) bool {
		return !fn(value)
	}
}

// Filter returns a new slice containing only elements that match the predicate.
// This function does not modify the input slice.
// The result has its capacity clipped to its length using slices.Clip.
// Returns nil if the input is empty or no elements match the predicate.
func Filter[T any](in []T, fn func(T) bool) []T {
	if len(in) == 0 {
		return nil
	}
	// Pre-allocate with input capacity for best performance
	result := make([]T, 0, len(in))
	for _, item := range in {
		if fn(item) {
			result = append(result, item)
		}
	}
	// Return nil if no items matched to be consistent with Map behavior
	if len(result) == 0 {
		return nil
	}
	// Clip to remove excess capacity
	return slices.Clip(result)
}

// FilterInPlace performs in-place filtering of a slice based on a predicate function.
// This function modifies the input slice and returns a sub-slice of the original.
// The returned slice shares the same backing array as the input.
// The result has its capacity clipped to its length using slices.Clip.
// Elements beyond the returned slice's length should not be accessed.
//
// This is more efficient than Filter when:
// - The input slice can be modified
// - You don't need to preserve the original data
// - You want to minimize allocations
//
// Example:
//
//	data := []int{1, 2, 3, 4, 5}
//	filtered := fp.FilterInPlace(data, func(n int) bool { return n%2 == 0 })
//	// filtered is []int{2, 4} backed by the same array as data
//	// data should not be used after this point
//
// Returns nil if the input is empty or no elements match the predicate.
func FilterInPlace[T any](in []T, fn func(T) bool) []T {
	if len(in) == 0 {
		return nil
	}

	// Optimized in-place filtering using the two-pointer technique
	n := 0
	for _, item := range in {
		if fn(item) {
			in[n] = item
			n++
		}
	}

	if n == 0 {
		return nil
	}

	// Clip to remove excess capacity
	return slices.Clip(in[:n])
}
