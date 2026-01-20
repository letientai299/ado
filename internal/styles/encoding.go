package styles

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-json"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
	"github.com/muesli/termenv"
)

var (
	yamlPrinter     *printer.Printer
	jsonColorScheme *json.ColorScheme
)

func YAML(v any) string {
	bs, err := encodeYAML(v)
	if err != nil {
		return ""
	}

	// Fast path for non-colored output
	if yamlPrinter == nil {
		return string(bs)
	}

	// Tokenize and colorize
	tokens := lexer.Tokenize(string(bs))
	return yamlPrinter.PrintTokens(tokens)
}

func DumpYAML(v any) error {
	bs, err := encodeYAML(v)
	if err != nil {
		log.Errorf("fail to marshal yaml: %v, err=%v", v, err)
		return err
	}

	// Fast path for non-colored output
	if yamlPrinter == nil {
		_, err = fmt.Print(string(bs))
		return err
	}

	// Tokenize and colorize
	tokens := lexer.Tokenize(string(bs))
	_, err = fmt.Print(yamlPrinter.PrintTokens(tokens))
	return err
}

func encodeYAML(v any) ([]byte, error) {
	bs, err := yaml.MarshalWithOptions(
		v,
		yaml.Indent(2),
		yaml.OmitEmpty(),
		yaml.OmitZero(),
		yaml.UseLiteralStyleIfMultiline(true),
	)
	return bs, err
}

// initYAMLPrinter creates the YAML printer with theme colors
// This is called once from Init() when the theme is set
func initYAMLPrinter(out *termenv.Output) {
	// Helper to create a property with ANSI escape codes
	printFn := func(color string) func() *printer.Property {
		if color == "" {
			return func() *printer.Property { return &printer.Property{} }
		}
		prefix, suffix := colorCodes(out, color)
		prop := &printer.Property{Prefix: prefix, Suffix: suffix}
		return func() *printer.Property { return prop }
	}

	// Create a printer with color properties
	yamlPrinter = &printer.Printer{
		Bool:    printFn(theme.Tokens.Chroma.Keyword),
		Number:  printFn(theme.Tokens.Chroma.Number),
		MapKey:  printFn(theme.Tokens.Chroma.NameTag),
		Anchor:  printFn(theme.Tokens.Chroma.NameConstant),
		Alias:   printFn(theme.Tokens.Chroma.NameConstant),
		String:  printFn(theme.Tokens.Chroma.String),
		Comment: printFn(theme.Tokens.Chroma.Comment),
	}
}

// initJSONColorScheme creates the JSON color scheme with theme colors
// This is called once from Init() when the theme is set
func initJSONColorScheme(out *termenv.Output) {
	// Helper to create a ColorFormat with ANSI escape codes
	colorFmt := func(color string) json.ColorFormat {
		if color == "" {
			return json.ColorFormat{}
		}
		prefix, suffix := colorCodes(out, color)
		return json.ColorFormat{Header: prefix, Footer: suffix}
	}

	jsonColorScheme = &json.ColorScheme{
		Int:       colorFmt(theme.Tokens.Chroma.Number),
		Uint:      colorFmt(theme.Tokens.Chroma.Number),
		Float:     colorFmt(theme.Tokens.Chroma.Number),
		Bool:      colorFmt(theme.Tokens.Chroma.Keyword),
		String:    colorFmt(theme.Tokens.Chroma.String),
		Binary:    colorFmt(theme.Tokens.Chroma.String),
		ObjectKey: colorFmt(theme.Tokens.Chroma.NameTag),
		Null:      colorFmt(theme.Tokens.Chroma.NameConstant),
	}
}

// Helper to create a property with ANSI escape codes
func colorCodes(out *termenv.Output, color string) (prefix, suffix string) {
	if color == "" {
		return "", ""
	}

	// Extract ANSI escape sequences
	const marker = "###"
	styled := out.String(marker).Foreground(out.Color(color)).Bold().String()
	if idx := strings.Index(styled, marker); idx >= 0 {
		prefix = styled[:idx]
		suffix = styled[idx+len(marker):]
	}
	return prefix, suffix
}

// DumpJSON prints the object as prettified JSON to stdout
func DumpJSON(v any) error {
	encoder := json.NewEncoder(os.Stdout)
	return encodeJSON(v, encoder)
}

func JSON(v any) string {
	var buf bytes.Buffer
	_ = encodeJSON(v, json.NewEncoder(&buf))
	return buf.String()
}

func encodeJSON(v any, encoder *json.Encoder) error {
	encoder.SetIndent("", "  ")

	// Use the theme color scheme if colors are enabled
	if jsonColorScheme != nil {
		err := encoder.EncodeWithOption(v, json.Colorize(jsonColorScheme))
		if err != nil {
			log.Errorf("fail to dump json: %v, err=%v", v, err)
			return err
		}
		return nil
	}

	err := encoder.Encode(v)
	if err != nil {
		log.Errorf("fail to dump json: %v, err=%v", v, err)
		return err
	}

	return nil
}
