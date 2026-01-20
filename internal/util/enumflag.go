package util

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/util/fp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ pflag.Value = (*EnumFlag[string])(nil)

// EnumFlag is a string flag that validates against a set of allowed values.
type EnumFlag[T ~string] struct {
	value        T   // flag value
	values       []T // all valid enum values
	valueStrings []string
	optional     bool
}

// NewEnumFlag creates an EnumFlag with the given default value and built-in allowed values.
func NewEnumFlag[T ~string](vals ...T) *EnumFlag[T] {
	return (&EnumFlag[T]{values: vals}).update()
}

func (e *EnumFlag[T]) Default(v T) *EnumFlag[T] {
	e.value = v
	return e
}

func (e *EnumFlag[T]) Optional() *EnumFlag[T] {
	e.optional = true
	return e
}

func (e *EnumFlag[T]) String() string { return string(e.value) }

func (e *EnumFlag[T]) Set(s string) error {
	// Don't validate here - custom values may be added later via AddAllowed.
	// Validation should be done after config resolution via Validate().
	e.value = T(s)
	return nil
}

// Validate checks if the current value is one of the allowed values.
// Call this after all values have been added via AddAllowed.
func (e *EnumFlag[T]) Validate() error {
	if e.optional && e.value == "" {
		return nil
	}

	for _, v := range e.values {
		if strings.EqualFold(string(e.value), string(v)) {
			e.value = v // normalize a case
			return nil
		}
	}
	return fmt.Errorf("must be one of: %s", strings.Join(e.valueStrings, ", "))
}

func (e *EnumFlag[T]) Type() string {
	return strings.Join(e.valueStrings, "|")
}

// AddAllowed adds additional allowed values (for runtime extension from config).
func (e *EnumFlag[T]) AddAllowed(vals ...T) {
	for _, v := range vals {
		if !slices.Contains(e.values, v) {
			e.values = append(e.values, v)
		}
	}
	e.update()
}

func (e *EnumFlag[T]) update() *EnumFlag[T] {
	slices.Sort(e.values)
	e.valueStrings = fp.Map(e.values, func(t T) string { return string(t) })
	return e
}

// Value returns the current value.
func (e *EnumFlag[T]) Value() T {
	if e == nil {
		return T("")
	}
	return e.value
}

// RegisterCompletion registers shell completion for this flag with cobra.
func (e *EnumFlag[T]) RegisterCompletion(cmd *cobra.Command, flagName string) {
	err := cmd.RegisterFlagCompletionFunc(
		flagName,
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return e.valueStrings, cobra.ShellCompDirectiveNoFileComp
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
