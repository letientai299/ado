package styles

import (
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/term"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
)

func ptr[T any](v T) *T { return &v }

const (
	defaultListIndent = 0
	defaultMargin     = 0
	MaxLineLength     = 100
)

var (
	mdRenderer *glamour.TermRenderer
	theme      = NoTTy
	out        = termenv.DefaultOutput()
)

func Init(th Theme) {
	if useColor() {
		theme = th
	}

	prepare()
}

func prepare() {
	initMdRenderer(theme)

	if theme.TrueColor {
		out = termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor))
	}

	initYAMLPrinter(out)
	initJSONColorScheme(out)
}

func useColor() bool {
	v := os.Getenv("COLOR")
	if v == "never" {
		return false
	}

	return v == "always" ||
		(term.IsTerminal(os.Stdout.Fd()) && term.IsTerminal(os.Stderr.Fd()))
}

func initMdRenderer(theme Theme) {
	var err error
	mdRenderer, err = glamour.NewTermRenderer(
		glamour.WithStyles(theme.glamourStyle()),
		glamour.WithWordWrap(MaxLineLength),
		glamour.WithInlineTableLinks(true),
	)
	if err != nil {
		log.Fatalf("fail to create markdown renderer: %v", err)
	}
}

func Heading(s string) string       { return colorize(s, theme.Tokens.Markdown.Heading) }
func Person(s string) string        { return colorize(s, theme.Tokens.Chroma.Name) }
func Time(s string) string          { return colorize(s, theme.Tokens.Chroma.Number) }
func Cmd(s string) string           { return colorize(s, theme.Tokens.Chroma.Function) }
func FlagStyle(s string) string     { return colorize(s, theme.Tokens.Chroma.Operator) }
func FlagTypeStyle(s string) string { return colorize(s, theme.Tokens.Chroma.KeywordType) }
func Warn(s string) string          { return colorize(s, theme.Tokens.Warn) }
func Success(s string) string       { return colorize(s, theme.Tokens.Success) }
func Error(s string) string         { return colorize(s, theme.Tokens.Error) }

func colorize(s, c string) string {
	return out.String(s).Foreground(out.Color(c)).Bold().String()
}

func Markdown(md string) string {
	s, err := mdRenderer.Render(md)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func Wrap(s string) string {
	w, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil || w <= 0 {
		w = MaxLineLength
	}
	if w > MaxLineLength {
		w = MaxLineLength
	}
	return wordwrap.String(s, w)
}
