package commands

import (
	"fmt"
	"strings"

	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func helpFunc(defaultHelp func(*cobra.Command, []string)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		_ = initConfig(cmd.Root())
		defaultHelp(cmd, args)
	}
}

func addTemplateHelpers() {
	cobra.AddTemplateFunc("headingStyle", styles.HeadingStyle)
	cobra.AddTemplateFunc("cmdStyle", styles.CmdStyle)
	cobra.AddTemplateFunc("renderFlags", renderFlags)
	cobra.AddTemplateFunc("markdown", styles.Markdown)
}

func renderFlags(cmd *cobra.Command) string {
	if !cmd.HasAvailableLocalFlags() && !cmd.HasAvailableInheritedFlags() {
		return ""
	}

	maxFlag := calcMaxFlagLen(cmd.LocalFlags(), cmd.InheritedFlags())
	var sb strings.Builder

	renderSection := func(title string, fs *pflag.FlagSet) {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
		sb.WriteString(styles.HeadingStyle(title))
		sb.WriteString("\n")
		sb.WriteString(renderFlagsWithMax(fs, maxFlag))
		sb.WriteString("\n")
	}

	if cmd.HasAvailableLocalFlags() {
		renderSection("Flags:", cmd.LocalFlags())
	}

	if cmd.HasAvailableInheritedFlags() {
		renderSection("Global Flags:", cmd.InheritedFlags())
	}

	return sb.String()
}

func calcMaxFlagLen(fss ...*pflag.FlagSet) flagMaxLen {
	maxFlag := flagMaxLen{}
	for _, fs := range fss {
		fs.VisitAll(func(f *pflag.Flag) {
			if f.Hidden {
				return
			}
			if n := len(flagNameOnly(f)); n > maxFlag.Name {
				maxFlag.Name = n
			}
			if t := len(flagType(f)); t > maxFlag.Type {
				maxFlag.Type = t
			}
		})
	}
	return maxFlag
}

func renderFlagsWithMax(fs *pflag.FlagSet, max flagMaxLen) string {
	var list []*pflag.Flag
	fs.VisitAll(func(f *pflag.Flag) {
		if !f.Hidden {
			list = append(list, f)
		}
	})

	if len(list) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, f := range list {
		sb.WriteString(formatFlag(f, max.Name, max.Type))
		if i < len(list)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func formatFlag(f *pflag.Flag, nameMax, typeMax int) string {
	name := flagNameOnly(f)
	fType := flagType(f)

	usage := f.Usage
	if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" && f.DefValue != "[]" &&
		f.DefValue != "\"\"" && f.DefValue != "<nil>" {
		usage += fmt.Sprintf(" (default %q)", f.DefValue)
	}

	// Calculate prefixes for alignment
	namePadding := strings.Repeat(" ", nameMax-len(name))
	typePadding := strings.Repeat(" ", typeMax-len(fType))

	// If there is a type, add a space before it
	typeStr := ""
	if fType != "" {
		typeStr = " " + styles.FlagTypeStyle(fType)
		// if we have type, the padding for type comes after it
		// but if we don't have the type, we still need to pad to typeMax + 1 (for the space)
	} else {
		typePadding = strings.Repeat(
			" ",
			typeMax+1,
		) // +1 for the space that would have been before the type
	}

	fullFlag := fmt.Sprintf("  %s%s%s%s", styles.FlagStyle(name), namePadding, typeStr, typePadding)

	// We need the unstyled length of fullFlag for calculation
	fullFlagLen := 2 + len(name) + len(namePadding) + 1 + len(fType) + len(typePadding)
	if fType == "" {
		fullFlagLen = 2 + len(name) + len(namePadding) + len(typePadding)
	}

	if fullFlagLen+len(usage)+1 <= styles.MaxLineLength {
		return fmt.Sprintf("%s %s", fullFlag, usage)
	}
	return fmt.Sprintf("%s\n%s", fullFlag, util.Indent(4, styles.Wrap(usage)))
}

type flagMaxLen struct {
	Name int
	Type int
}

func flagNameOnly(f *pflag.Flag) string {
	if f.Shorthand != "" {
		return "-" + f.Shorthand + ", --" + f.Name
	}
	return "    --" + f.Name
}

func flagType(f *pflag.Flag) string {
	if f.Value.Type() == "bool" {
		return ""
	}
	return f.Value.Type()
}
