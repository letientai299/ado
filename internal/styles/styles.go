package styles

import (
	"os"
	"sync"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/term"
	"github.com/muesli/termenv"
)

func ptr[T any](v T) *T { return &v }

const (
	defaultListIndent = 0
	defaultMargin     = 0
	MaxLineLength     = 80
)

var (
	mdRenderer     *glamour.TermRenderer
	mdRendererOnce sync.Once
	theme          = NoTTy
	out            = termenv.DefaultOutput()
	useColor       = computeUseColor()
)

func Init(th Theme) {
	if UseColor() {
		theme = th
	}

	prepare()
}

func prepare() {
	useColor = computeUseColor()
	switch {
	case !useColor:
		out = termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.Ascii))
		lipgloss.SetColorProfile(termenv.Ascii)
	case theme.TrueColor:
		out = termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor))
	default:
		out = termenv.DefaultOutput()
	}

	precomputeStyleCache(theme)
	initYAMLPrinter()
	initJSONColorScheme()
}

func UseColor() bool {
	return useColor
}

func computeUseColor() bool {
	if termenv.EnvNoColor() {
		return false
	}

	v := os.Getenv("COLOR")
	if v == "never" {
		return false
	}

	if v == "always" {
		return true
	}

	return term.IsTerminal(os.Stdout.Fd()) && term.IsTerminal(os.Stderr.Fd())
}

// initMdRenderer lazily initializes the glamour markdown renderer.
// Note: chroma's lexer init happens at import time and cannot be deferred.
func initMdRenderer() {
	mdRendererOnce.Do(func() {
		var err error
		mdRenderer, err = glamour.NewTermRenderer(
			glamour.WithStyles(theme.glamourStyle()),
			glamour.WithWordWrap(MaxLineLength),
			glamour.WithInlineTableLinks(true),
		)
		if err != nil {
			log.Fatalf("fail to create markdown renderer: %v", err)
		}
	})
}

func Heading(s string) string       { return colorize(s, theme.Tokens.Markdown.Heading) }
func H1(s string) string            { return colorize(s, theme.Tokens.Markdown.H1) }
func Person(s string) string        { return colorize(s, theme.Tokens.Chroma.Name) }
func Time(s string) string          { return colorize(s, theme.Tokens.Chroma.Number) }
func Const(s string) string         { return colorize(s, theme.Tokens.Chroma.NameConstant) }
func Cmd(s string) string           { return colorize(s, theme.Tokens.Chroma.Function) }
func FlagStyle(s string) string     { return colorize(s, theme.Tokens.Chroma.Operator) }
func FlagTypeStyle(s string) string { return colorize(s, theme.Tokens.Chroma.KeywordType) }
func Warn(s string) string          { return colorize(s, theme.Tokens.Warn) }
func Success(s string) string       { return colorize(s, theme.Tokens.Success) }
func Pending(s string) string       { return colorize(s, theme.Tokens.Chroma.Function) }
func Highlight(s string) string     { return underline(colorize(s, theme.Tokens.Success)) }
func Faint(s string) string         { return out.String(s).Faint().String() }
func Error(s string) string         { return colorize(s, theme.Tokens.Error) }

// Markdown renders content with syntax highlighting.
// The renderer is initialized lazily on the first call.
func Markdown(md string) string {
	initMdRenderer()
	s, err := mdRenderer.Render(md)
	if err != nil {
		log.Fatal(err)
	}
	return s
}
