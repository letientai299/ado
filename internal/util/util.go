package util

type StrErr string

func (s StrErr) Error() string { return string(s) }

func Ptr[T any](v T) *T { return &v }

func PanicIfNil(v any) {
	if v == nil {
		panic("nil value")
	}
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
