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

var (
	mdRenderer *glamour.TermRenderer
	theme      = NoTTy
	out        = termenv.DefaultOutput()
)

func Init(th Theme) {
	theme = th
	initMdRenderer(th)

	if th.TrueColor {
		out = termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor))
	}

	initYAMLPrinter(out, th)
	initJSONColorScheme(out, th)
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

func HeadingStyle(s string) string  { return colorize(s, theme.Tokens.Markdown.Heading) }
func CmdStyle(s string) string      { return colorize(s, theme.Tokens.Chroma.Function) }
func FlagStyle(s string) string     { return colorize(s, theme.Tokens.Chroma.Operator) }
func FlagTypeStyle(s string) string { return colorize(s, theme.Tokens.Chroma.KeywordType) }

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
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		w = MaxLineLength
	}
	if w > MaxLineLength {
		w = MaxLineLength
	}
	return wordwrap.String(s, w)
}
