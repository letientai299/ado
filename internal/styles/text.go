package styles

import (
	"bytes"
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
	if s == "" {
		return
	}

	padding := strings.Repeat(" ", n)
	iw := &indentWriter{w: w, indent: []byte(padding), atBOL: true}
	_, err := iw.Write([]byte(s))
	util.PanicIf(err)

	if iw.atBOL {
		_, err = w.Write(iw.indent)
		util.PanicIf(err)
	}
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
	indent []byte
	atBOL  bool // at beginning of line
}

func NewIndentWriter(w io.Writer, indent string) io.Writer {
	return &indentWriter{w: w, indent: []byte(indent)}
}

func (iw *indentWriter) Write(bs []byte) (n int, err error) {
	for len(bs) > 0 {
		if iw.atBOL {
			if bs[0] == '\n' {
				// Multiple newlines: write the \n and stay at BOL
				nn, err := iw.w.Write(bs[:1])
				n += nn
				if err != nil {
					return n, err
				}
				bs = bs[1:]
				continue
			}

			if _, err := iw.w.Write(iw.indent); err != nil {
				return n, err
			}
			iw.atBOL = false
		}

		idx := bytes.IndexByte(bs, '\n')
		if idx == -1 {
			nn, err := iw.w.Write(bs)
			n += nn
			if err != nil {
				return n, err
			}
			return n, nil
		}

		nn, err := iw.w.Write(bs[:idx+1])
		n += nn
		if err != nil {
			return n, err
		}
		iw.atBOL = true
		bs = bs[idx+1:]
	}

	return n, nil
}
