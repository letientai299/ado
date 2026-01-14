package _shared

import "io"

type Querier interface {
	AppendTo(io.Writer)
}

var (
	_ Querier = BoolQ("")
	_ Querier = Queriers{}
)

type BoolQ string

func (b BoolQ) AppendTo(w io.Writer) {
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
