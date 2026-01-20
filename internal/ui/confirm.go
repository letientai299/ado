package ui

import (
	"fmt"
	"strings"

	"github.com/letientai299/ado/internal/styles"
)

func Confirm(message string, defaultValue bool) bool {
	var choices string
	if defaultValue {
		choices = " [" + styles.Success("Y") + "/n]"
	} else {
		choices = " [y/" + styles.Success("N") + "]"
	}
	fmt.Printf("%s%s: ", styles.Warn("! "), message+choices)
	var input string
	_, _ = fmt.Scanln(&input)
	input = strings.ToLower(strings.TrimSpace(input))
	if input == "" {
		return defaultValue
	}
	return input == "y" || input == "yes"
}
