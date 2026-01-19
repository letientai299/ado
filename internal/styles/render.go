package styles

import (
	"io"
	"os"
	"strings"
	"text/template"
)

var TemplateFuncs = template.FuncMap{
	"warn":      Warn,
	"error":     Error,
	"success":   Success,
	"heading":   Heading,
	"h1":        H1,
	"person":    Person,
	"time":      Time,
	"cmdStyle":  Cmd,
	"markdown":  Markdown,
	"join":      strings.Join,
	"indent":    Indent,
	"trimSpace": strings.TrimSpace,
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

func RenderOut(tpl string, v any, funcMaps ...template.FuncMap) error {
	return Render(os.Stdout, tpl, v, funcMaps...)
}
