package styles

import (
	"io"
	"os"
	"strings"
	"text/template"
)

var TemplateFuncs = template.FuncMap{
	"const":      Const,
	"faint":      Faint,
	"warn":       Warn,
	"error":      Error,
	"success":    Success,
	"pending":    Pending,
	"highlight":  Highlight,
	"heading":    Heading,
	"h1":         H1,
	"person":     Person,
	"time":       Time,
	"cmdStyle":   Cmd,
	"markdown":   Markdown,
	"join":       func(sep string, s []string) string { return strings.Join(s, sep) },
	"indent":     Indent,
	"trimSpace":  strings.TrimSpace,
	"trimLeft":   func(cutset, s string) string { return strings.TrimLeft(s, cutset) },
	"trimPrefix": func(prefix, s string) string { return strings.TrimPrefix(s, prefix) },
	"replaceAll": func(old, new, s string) string { return strings.ReplaceAll(s, old, new) },
}

func Render(w io.Writer, tpl string, v any, funcMaps ...template.FuncMap) error {
	t := template.New("output")
	t.Funcs(TemplateFuncs)
	for _, m := range funcMaps {
		t.Funcs(m)
	}

	t, err := t.Parse(tpl)
	if err != nil {
		return err
	}

	return t.Execute(w, v)
}

func RenderS(tpl string, v any, funcMaps ...template.FuncMap) (string, error) {
	var sb strings.Builder
	err := Render(&sb, tpl, v, funcMaps...)
	return sb.String(), err
}

func RenderOut(tpl string, v any, funcMaps ...template.FuncMap) error {
	return Render(os.Stdout, tpl, v, funcMaps...)
}

func HighlightMatch(s string, matches []int) string {
	if len(matches) == 0 {
		return s
	}

	var b strings.Builder
	runes := []rune(s)
	matchSet := make(map[int]bool)
	for _, m := range matches {
		matchSet[m] = true
	}

	for i, r := range runes {
		if matchSet[i] {
			b.WriteString(Highlight(string(r)))
		} else {
			b.WriteRune(r)
		}
	}

	return b.String()
}
