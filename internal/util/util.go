package util

type StrErr string

func (s StrErr) Error() string { return string(s) }

func Ptr[T any](v T) *T { return &v }
