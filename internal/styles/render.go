package styles

import (
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
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
	"trimPrefix": func(prefix, s string) string { return strings.TrimPrefix(s, prefix) },
	"replaceAll": func(old, new, s string) string { return strings.ReplaceAll(s, old, new) },
	"tr":         func(old, new, s string) string { return regexp.MustCompile(old).ReplaceAllString(s, new) },
}

// templateCache caches parsed templates to avoid repeated parsing.
// Only cache templates without custom funcMaps (the common case).
var (
	templateCache   = make(map[string]*template.Template)
	templateCacheMu sync.RWMutex
)

// baseTemplate is a pre-configured template with standard funcs, used as a clone source.
var baseTemplate = template.New("base").Funcs(TemplateFuncs)

func Render(w io.Writer, tpl string, v any, funcMaps ...template.FuncMap) error {
	// Fast path: no custom funcMaps, use cached template
	if len(funcMaps) == 0 {
		t, err := getCachedTemplate(tpl)
		if err != nil {
			return err
		}
		return t.Execute(w, v)
	}

	// Slow path: custom funcMaps, create new template
	t, err := baseTemplate.Clone()
	if err != nil {
		return err
	}
	for _, m := range funcMaps {
		t.Funcs(m)
	}
	t, err = t.Parse(tpl)
	if err != nil {
		return err
	}
	return t.Execute(w, v)
}

// getCachedTemplate returns a cached parsed template, parsing it on first access.
func getCachedTemplate(tpl string) (*template.Template, error) {
	// Try read lock first (fast path)
	templateCacheMu.RLock()
	if t, ok := templateCache[tpl]; ok {
		templateCacheMu.RUnlock()
		return t, nil
	}
	templateCacheMu.RUnlock()

	// Parse and cache with write lock
	templateCacheMu.Lock()
	defer templateCacheMu.Unlock()

	// Double-check after acquiring write lock
	if t, ok := templateCache[tpl]; ok {
		return t, nil
	}

	t, err := baseTemplate.Clone()
	if err != nil {
		return nil, err
	}
	t, err = t.Parse(tpl)
	if err != nil {
		return nil, err
	}
	templateCache[tpl] = t
	return t, nil
}

// ClearTemplateCache clears the template cache. Useful for testing.
func ClearTemplateCache() {
	templateCacheMu.Lock()
	defer templateCacheMu.Unlock()
	templateCache = make(map[string]*template.Template)
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
