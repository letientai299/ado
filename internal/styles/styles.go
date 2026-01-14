package styles

import (
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/log"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

func ptr[T any](v T) *T { return &v }

const (
	defaultListIndent = 0
	defaultMargin     = 0
	MaxLineLength     = 100
)

var mdRenderer *glamour.TermRenderer

var (
	HeadingStyle func(string) string
	FlagStyle    func(string) string
	CmdStyle     func(string) string
)

func Init(theme Theme) {
	initMdRenderer(theme)

	out := termenv.DefaultOutput()
	if theme.TrueColor {
		out = termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor))
	}

	initColorizers(out, theme)
	initYAMLPrinter(out, theme)
	initJSONColorScheme(out, theme)
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

func initColorizers(out *termenv.Output, theme Theme) {
	HeadingStyle = colorize(out, theme.Tokens.Heading)
	CmdStyle = colorize(out, theme.Tokens.Chroma.Function)
	FlagStyle = colorize(out, theme.Tokens.Chroma.Operator)
}

func colorize(out *termenv.Output, c string) func(string) string {
	return func(s string) string {
		return out.String(s).Foreground(out.Color(c)).Bold().String()
	}
}

func Markdown(md string) (string, error) {
	return mdRenderer.Render(md)
}

func Wrap(s string) string {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		w = MaxLineLength
	}
	if w > MaxLineLength {
		w = MaxLineLength
	}
	return wordwrap.String(s, w)
}
