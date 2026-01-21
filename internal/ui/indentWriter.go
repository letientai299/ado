package ui

import "io"

var _ io.Writer = &indentWriter{}

// indentWriter add indent to written content for every newline.
type indentWriter struct {
	w      io.Writer
	indent string
	atBOL  bool // at beginning of line
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

func NewIndentWriter(w io.Writer, indent string) io.Writer {
	return &indentWriter{w: w, indent: indent}
}
