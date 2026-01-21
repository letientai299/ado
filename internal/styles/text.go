package styles

import (
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/x/term"
	"github.com/letientai299/ado/internal/util"
	"github.com/muesli/reflow/wordwrap"
)

// Indent add indentation of n spaces to every line in the string
func Indent(n int, s string) string {
	var sb strings.Builder
	IndentTo(&sb, n, s)
	return sb.String()
}

func IndentTo(w io.Writer, n int, s string) {
	padding := strings.Repeat(" ", n)
	_, err := w.Write([]byte(padding))
	util.PanicIf(err)

	_, err = w.Write([]byte(strings.ReplaceAll(s, "\n", "\n"+padding)))
	util.PanicIf(err)
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

var _ io.Writer = &indentWriter{}

// indentWriter add indent to written content for every newline.
type indentWriter struct {
	w      io.Writer
	indent string
	atBOL  bool // at beginning of line
}

func NewIndentWriter(w io.Writer, indent string) io.Writer {
	return &indentWriter{w: w, indent: indent}
}

func (iw *indentWriter) Write(bs []byte) (n int, err error) {
	for _, b := range bs {
		if iw.atBOL && b != '\n' {
			_, err = iw.w.Write([]byte(iw.indent))
			if err != nil {
				return n, err
			}
			iw.atBOL = false
		}

		_, err = iw.w.Write([]byte{b})
		if err != nil {
			return n, err
		}
		n++
		if b == '\n' {
			iw.atBOL = true
		}
	}

	return n, nil
}
