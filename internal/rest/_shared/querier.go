package _shared

import (
	"fmt"
	"io"
)

type Querier interface {
	AppendTo(io.Writer)
}

var (
	_ Querier = Bool("")
	_ Querier = Queriers{}
)

type Bool string

func (b Bool) AppendTo(w io.Writer) {
	if b != "" {
		_, _ = w.Write([]byte("&" + string(b) + "=true"))
	}
}

type Queriers []Querier

func (qs Queriers) AppendTo(w io.Writer) {
	for _, q := range qs {
		q.AppendTo(w)
	}
}

type KV[T any] struct {
	Key   string
	Value T
}

func (kv KV[T]) AppendTo(w io.Writer) {
	_, _ = w.Write([]byte("&" + kv.Key + "="))
	_, _ = fmt.Fprint(w, kv.Value)
}
