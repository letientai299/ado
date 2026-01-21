package styles

import (
	"strings"
	"unicode/utf8"
)

var styleCache = make(map[string]ansiStyle)

type ansiStyle struct {
	prefix string
	suffix string
}

func getStyle(color string) (string, string) {
	if st, ok := styleCache[color]; ok {
		return st.prefix, st.suffix
	}

	return computeStyle(color)
}

func computeStyle(color string) (string, string) {
	if color == "" {
		return "", ""
	}

	const dummy = "@"
	full := out.String(dummy).Foreground(out.Color(color)).Bold().String()
	idx := strings.Index(full, dummy)
	if idx == -1 {
		return "", ""
	}
	st := ansiStyle{
		prefix: full[:idx],
		suffix: full[idx+len(dummy):],
	}
	return st.prefix, st.suffix
}

func clearStyleCache() {
	styleCache = make(map[string]ansiStyle)
}

func precomputeStyleCache(th Theme) {
	clearStyleCache()

	colors := []string{
		th.Tokens.Warn,
		th.Tokens.Error,
		th.Tokens.Success,
		th.Tokens.Markdown.Text,
		th.Tokens.Markdown.Heading,
		th.Tokens.Markdown.H1,
		th.Tokens.Markdown.H1Background,
		th.Tokens.Markdown.HorizontalRule,
		th.Tokens.Markdown.Link,
		th.Tokens.Markdown.LinkText,
		th.Tokens.Markdown.Image,
		th.Tokens.Markdown.ImageText,
		th.Tokens.Markdown.Code,
		th.Tokens.Markdown.CodeBackground,
		th.Tokens.Markdown.CodeBlock,
		th.Tokens.Chroma.Text,
		th.Tokens.Chroma.Error,
		th.Tokens.Chroma.ErrorBackground,
		th.Tokens.Chroma.Comment,
		th.Tokens.Chroma.CommentPreproc,
		th.Tokens.Chroma.Keyword,
		th.Tokens.Chroma.KeywordReserved,
		th.Tokens.Chroma.KeywordNamespace,
		th.Tokens.Chroma.KeywordType,
		th.Tokens.Chroma.Operator,
		th.Tokens.Chroma.Punctuation,
		th.Tokens.Chroma.Name,
		th.Tokens.Chroma.NameBuiltin,
		th.Tokens.Chroma.NameTag,
		th.Tokens.Chroma.NameAttribute,
		th.Tokens.Chroma.NameClass,
		th.Tokens.Chroma.NameConstant,
		th.Tokens.Chroma.NameDecorator,
		th.Tokens.Chroma.Function,
		th.Tokens.Chroma.String,
		th.Tokens.Chroma.LiteralStringEscape,
		th.Tokens.Chroma.Number,
		th.Tokens.Chroma.GenericDeleted,
		th.Tokens.Chroma.GenericInserted,
		th.Tokens.Chroma.GenericSubheading,
		th.Tokens.Chroma.Background,
	}

	for _, c := range colors {
		if c == "" {
			continue
		}
		if _, ok := styleCache[c]; !ok {
			p, s := computeStyle(c)
			styleCache[c] = ansiStyle{prefix: p, suffix: s}
		}
	}
}

func applyStyle(s, color string) string {
	if s == "" {
		return s
	}

	prefix, suffix := getStyle(color)
	return prefix + s + suffix
}

func colorize(s, color string) string {
	if s == "" || !UseColor() {
		return s
	}

	// Fast path: no ANSI escape codes
	if strings.IndexByte(s, '\x1b') == -1 {
		return applyStyle(s, color)
	}

	prefix, suffix := getStyle(color)

	var b strings.Builder
	b.Grow(len(s) + 64) // Estimate extra space for ANSI codes
	last := 0
	styled := false

	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			// Found ANSI escape sequence
			text := s[last:i]
			if text != "" {
				if styled {
					b.WriteString(text)
				} else {
					b.WriteString(prefix)
					b.WriteString(text)
					b.WriteString(suffix)
				}
			}

			// Find end of ANSI sequence
			start := i
			i += 2 // skip \x1b[
			for i < len(s) {
				c := s[i]
				i++
				if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
					break
				}
			}
			ansi := s[start:i]
			b.WriteString(ansi)
			if strings.HasSuffix(ansi, "0m") {
				styled = false
			} else {
				styled = true
			}
			last = i
		} else {
			_, size := utf8.DecodeRuneInString(s[i:])
			i += size
		}
	}

	text := s[last:]
	if text != "" {
		if styled {
			b.WriteString(text)
		} else {
			b.WriteString(prefix)
			b.WriteString(text)
			b.WriteString(suffix)
		}
	}
	return b.String()
}

func underline(s string) string {
	if s == "" || !UseColor() {
		return s
	}
	return out.String(s).Underline().String()
}
