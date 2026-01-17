package util

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ pflag.Value = (*EnumFlag)(nil)

// EnumFlag is a string flag that validates against a set of allowed values.
type EnumFlag struct {
	value  string   // flag value
	values []string // all valid enum values
}

// NewEnumFlag creates an EnumFlag with the given default value and built-in allowed values.
func NewEnumFlag(defValue string, vals ...string) *EnumFlag {
	sort.Strings(vals)
	return &EnumFlag{
		value:  defValue,
		values: vals,
	}
}

func (e *EnumFlag) String() string { return e.value }

func (e *EnumFlag) Set(s string) error {
	// Don't validate here - custom values may be added later via AddAllowed.
	// Validation should be done after config resolution via Validate().
	e.value = s
	return nil
}

// Validate checks if the current value is one of the allowed values.
// Call this after all values have been added via AddAllowed.
func (e *EnumFlag) Validate() error {
	for _, v := range e.values {
		if strings.EqualFold(e.value, v) {
			e.value = v // normalize case
			return nil
		}
	}
	return fmt.Errorf("must be one of: %s", strings.Join(e.values, ", "))
}

func (e *EnumFlag) Type() string {
	return strings.Join(e.values, "|")
}

// AddAllowed adds additional allowed values (for runtime extension from config).
func (e *EnumFlag) AddAllowed(vals ...string) {
	for _, v := range vals {
		if !slices.Contains(e.values, v) {
			e.values = append(e.values, v)
		}
	}
	sort.Strings(e.values)
}

// Value returns the current value.
func (e *EnumFlag) Value() string {
	return e.value
}

// RegisterCompletion registers shell completion for this flag with cobra.
func (e *EnumFlag) RegisterCompletion(cmd *cobra.Command, flagName string) error {
	return cmd.RegisterFlagCompletionFunc(
		flagName,
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return slices.Clone(e.values), cobra.ShellCompDirectiveNoFileComp
		},
	)
}
