package util

import (
	"github.com/spf13/cobra"
)

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

// HelpTopic creates a help topic command with the given use and documentation.
func HelpTopic(use, doc string) *cobra.Command {
	return &cobra.Command{Use: use, Long: doc}
}