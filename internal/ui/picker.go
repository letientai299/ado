package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/util"
	"github.com/letientai299/ado/internal/util/fp"
)

var (
	_ list.Item         = (*pickItem[any])(nil)
	_ list.ItemDelegate = (*pickDelegate[any])(nil)
)

const (
	errPickConfigNeedRender util.StrErr = "function Render is required"
	errPickConfigNeedFilter util.StrErr = "function FilterValue is required"
)

type PickConfig[T any] struct {
	Render      func(w io.Writer, it T, matches []int)
	FilterValue func(item T) string
	ItemHeight  int
}

func (p *PickConfig[T]) validate() error {
	if p.Render == nil {
		return errPickConfigNeedRender
	}

	if p.FilterValue == nil {
		return errPickConfigNeedFilter
	}

	if p.ItemHeight <= 0 {
		p.ItemHeight = 1
	}

	return nil
}

// Pick allows the user to pick an item from a list.
//
// It should support fuzzy search. User could start typing immediately when the list appears to
// filter items.
//
// It supports these key bindings: ctrl-n or up for "next", ctrl-p or down for "prev", ctrl-c
// cancel/quit, and enter for selection.
//
// The caller controls item style via the PickConfig.Render function.
// Other UI components will use the color tokens from styles.GetTheme.
func Pick[T any](items []T, cfg PickConfig[T]) fp.Optional[T] {
	if len(items) == 0 {
		return fp.Nil[T]()
	}

	err := cfg.validate()
	util.PanicIf(err)

	// Convert items to list.Item
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		// FilterValue is guaranteed to be valid after validation
		listItems[i] = pickItem[T]{
			value:       item,
			filterValue: cfg.FilterValue(item),
		}
	}

	delegate := &pickDelegate[T]{
		items: items,
		cfg:   cfg,
	}

	model := pickerModel[T]{
		list:   newList(listItems, delegate),
		cfg:    cfg,
		picked: fp.Nil[T](),
	}

	p := tea.NewProgram(model, tea.WithOutput(os.Stderr))
	finalModel, err := p.Run()
	if err != nil {
		return fp.Nil[T]()
	}

	return finalModel.(pickerModel[T]).picked
}

func newList[T any](listItems []list.Item, delegate *pickDelegate[T]) list.Model {
	l := list.New(listItems, delegate, 0, 0)
	noMargin := func(s *lipgloss.Style) {
		*s = s.Margin(0).PaddingTop(0).PaddingBottom(0).PaddingLeft(0)
	}
	for _, s := range []*lipgloss.Style{
		&l.Styles.Title,
		&l.Styles.FilterPrompt,
		&l.Styles.StatusBar,
		&l.Styles.StatusBarFilterCount,
		&l.Styles.StatusBarActiveFilter,
		&l.Styles.StatusEmpty,
		&l.Styles.NoItems,
		&l.Styles.PaginationStyle,
		&l.Styles.HelpStyle,
		&l.Styles.TitleBar,
	} {
		noMargin(s)
	}

	l.SetShowTitle(false)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)
	return l
}

// pickItem wraps a value of type T to implement list.Item interface
type pickItem[T any] struct {
	value       T
	filterValue string
}

func (pi pickItem[T]) FilterValue() string { return pi.filterValue }

type pickDelegate[T any] struct {
	items []T
	cfg   PickConfig[T]
}

func (pd pickDelegate[T]) Height() int                               { return pd.cfg.ItemHeight }
func (pd pickDelegate[T]) Spacing() int                              { return 0 }
func (pd pickDelegate[T]) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (pd pickDelegate[T]) Render(w io.Writer, m list.Model, filteredIndex int, it list.Item) {
	if m.Index() == filteredIndex {
		_, _ = fmt.Fprint(w, styles.Success("▶ "))
	} else {
		_, _ = fmt.Fprint(w, "  ")
	}

	pi := it.(pickItem[T])
	pd.cfg.Render(w, pi.value, m.MatchesForItem(filteredIndex))
}

// pickerModel is the tea.Model for the picker
type pickerModel[T any] struct {
	list   list.Model
	cfg    PickConfig[T]
	picked fp.Optional[T]
}

func (m pickerModel[T]) Init() tea.Cmd { return nil }

func (m pickerModel[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+n":
			m.list.CursorDown()
			return m, nil
		case "ctrl+p":
			m.list.CursorUp()
			return m, nil
		case "enter":
			i, ok := m.list.SelectedItem().(pickItem[T])
			if ok {
				m.picked = fp.Some(i.value)
				return m, tea.Quit
			}

		case "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h := msg.Height
		// extra lines for: filter, help
		maxH := len(m.list.Items())*m.cfg.ItemHeight + 2
		if h > maxH {
			h = maxH
		}
		m.list.SetSize(msg.Width, h)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m pickerModel[T]) View() string { return m.list.View() }
