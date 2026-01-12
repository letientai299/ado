package styles

import (
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"
	"github.com/charmbracelet/log"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

func boolPtr(b bool) *bool       { return &b }
func stringPtr(s string) *string { return &s }
func uintPtr(u uint) *uint       { return &u }

const (
	defaultListIndent = 0
	defaultMargin     = 0
	MaxLineLength     = 100
)

var UseColor bool

var (
	HeadingStyle func(string) string
	FlagStyle    func(string) string
	CmdStyle     func(string) string
	mdRenderer   *glamour.TermRenderer
)

func init() {
	UseColor = os.Getenv("COLOR") == "always" ||
		(term.IsTerminal(int(os.Stdout.Fd())) && term.IsTerminal(int(os.Stderr.Fd())))

	initMdRenderer()
	initUsageColorizers()
}

func initMdRenderer() {
	style := chooseBestStyle()
	var err error
	mdRenderer, err = glamour.NewTermRenderer(
		glamour.WithStyles(style),
		glamour.WithWordWrap(MaxLineLength),
		glamour.WithInlineTableLinks(true),
	)
	if err != nil {
		log.Fatalf("fail to create markdown renderer: %v", err)
	}
}

func initUsageColorizers() {
	var h, f, c string
	if termenv.HasDarkBackground() {
		h = *MdDark.Heading.Color
		f = *MdDark.Code.Color
		c = *MdDark.Code.Color
	} else {
		h = *MdLight.Heading.Color
		f = *MdLight.Code.Color
		c = *MdLight.Code.Color
	}

	HeadingStyle = colorize(h)
	FlagStyle = colorize(f)
	CmdStyle = colorize(c)
}

func colorize(c string) func(string) string {
	out := termenv.DefaultOutput()
	if os.Getenv("COLOR") == "always" {
		out = termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor))
	}
	return func(s string) string {
		return out.String(s).Foreground(out.Color(c)).Bold().String()
	}
}

func chooseBestStyle() ansi.StyleConfig {
	if !UseColor {
		return styles.NoTTYStyleConfig
	}

	if termenv.HasDarkBackground() {
		return MdDark
	}

	return MdLight
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
