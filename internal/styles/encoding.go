package styles

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-json"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
)

var (
	yamlPrinter     *printer.Printer
	jsonColorScheme *json.ColorScheme
	bufferPool      = sync.Pool{
		New: func() any {
			return new(bytes.Buffer)
		},
	}
)

func getBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func putBuffer(buf *bytes.Buffer) {
	bufferPool.Put(buf)
}

func YAML(v any) string {
	buf := getBuffer()
	defer putBuffer(buf)
	if err := encodeYAMLTo(buf, v); err != nil {
		return ""
	}

	// Fast path for non-colored output
	if !UseColor() || yamlPrinter == nil {
		return buf.String()
	}

	// Tokenize and colorize
	return yamlPrinter.PrintTokens(lexer.Tokenize(buf.String()))
}

func DumpYAML(v any) error {
	return DumpYAMLTo(os.Stdout, v)
}

func DumpYAMLTo(w io.Writer, v any) error {
	buf := getBuffer()
	defer putBuffer(buf)
	if err := encodeYAMLTo(buf, v); err != nil {
		log.Errorf("fail to marshal yaml: %v, err=%v", v, err)
		return err
	}

	// Fast path for non-colored output
	if !UseColor() || yamlPrinter == nil {
		_, err := w.Write(buf.Bytes())
		return err
	}

	// Tokenize and colorize
	_, err := fmt.Fprintln(w, yamlPrinter.PrintTokens(lexer.Tokenize(buf.String())))
	return err
}

func encodeYAMLTo(w io.Writer, v any) error {
	encoder := yaml.NewEncoder(
		w,
		yaml.Indent(2),
		yaml.OmitEmpty(),
		yaml.OmitZero(),
		yaml.UseLiteralStyleIfMultiline(true),
	)
	return encoder.Encode(v)
}

// initYAMLPrinter creates the YAML printer with theme colors
// This is called once from Init() when the theme is set
func initYAMLPrinter() {
	// Helper to create a property with ANSI escape codes
	printFn := func(color string) func() *printer.Property {
		if color == "" {
			return func() *printer.Property { return &printer.Property{} }
		}
		prefix, suffix := getStyle(color)
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
func initJSONColorScheme() {
	// Helper to create a ColorFormat with ANSI escape codes
	colorFmt := func(color string) json.ColorFormat {
		prefix, suffix := getStyle(color)
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

// DumpJSON prints the object as prettified JSON to writer
func DumpJSON(v any) error {
	return encodeJSON(v, json.NewEncoder(os.Stdout))
}

func DumpJSONTo(w io.Writer, v any) error {
	return encodeJSON(v, json.NewEncoder(w))
}

func JSON(v any) string {
	buf := getBuffer()
	defer putBuffer(buf)
	if err := encodeJSON(v, json.NewEncoder(buf)); err != nil {
		return ""
	}
	return buf.String()
}

func encodeJSON(v any, encoder *json.Encoder) error {
	encoder.SetIndent("", "  ")

	// Use the theme color scheme if colors are enabled
	if UseColor() && jsonColorScheme != nil {
		return encoder.EncodeWithOption(v, json.Colorize(jsonColorScheme))
	}

	return encoder.Encode(v)
}
