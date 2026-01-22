// Package fp provides functional programming utilities.
package fp

import "slices"

func Map[I, O any](in []I, fn func(I) O) []O {
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

func Filter[T any](in []T, fn func(T) bool) []T {
	return slices.DeleteFunc(slices.Clone(in), Not(fn))
}
