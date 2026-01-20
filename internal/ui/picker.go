package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
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
	Title       string
	Render      func(w io.Writer, it T, matches []int)
	FilterValue func(item T) string
	ItemHeight  int
}

func (pc *PickConfig[T]) toListItem(t T) list.Item {
	return pickItem[T]{
		value:       t,
		filterValue: pc.FilterValue(t),
	}
}

func (pc *PickConfig[T]) validate() error {
	if pc.Render == nil {
		return errPickConfigNeedRender
	}

	if pc.FilterValue == nil {
		return errPickConfigNeedFilter
	}

	if pc.ItemHeight <= 0 {
		pc.ItemHeight = 1
	}

	return nil
}

// Pick allows the user to pick an item from a list, supports fuzzy search and keyboard navigations.
func Pick[T any](items []T, cfg PickConfig[T]) fp.Optional[T] {
	if len(items) == 0 {
		return fp.Nil[T]()
	}

	util.PanicIf(cfg.validate())
	model := newPickModel(items, cfg)
	prog := tea.NewProgram(model, tea.WithOutput(os.Stderr))
	finalModel, err := prog.Run()
	if err != nil {
		log.Error("fail to pick item", "err", err)
		return fp.Nil[T]()
	}

	return finalModel.(*pickerModel[T]).picked
}

func newPickModel[T any](items []T, cfg PickConfig[T]) *pickerModel[T] {
	listItems := fp.Map(items, cfg.toListItem)
	return &pickerModel[T]{
		list:   newList(listItems, &pickDelegate[T]{cfg: cfg}),
		cfg:    cfg,
		picked: fp.Nil[T](),
	}
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

	if delegate.cfg.Title != "" {
		l.Title = delegate.cfg.Title
	} else {
		l.SetShowTitle(false)
	}
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
	cfg PickConfig[T]
}

func (pd pickDelegate[T]) Height() int                             { return pd.cfg.ItemHeight }
func (pd pickDelegate[T]) Spacing() int                            { return 0 }
func (pd pickDelegate[T]) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

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
	list     list.Model
	cfg      PickConfig[T]
	picked   fp.Optional[T]
	quitting bool

	maxHeight int // lazy computed
}

func (m *pickerModel[T]) Init() tea.Cmd {
	// extra lines for: filter, help
	m.maxHeight = len(m.list.Items())*m.cfg.ItemHeight + 2
	if m.cfg.Title != "" {
		m.maxHeight += 1
	}
	return nil
}

func (m *pickerModel[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.quitting = true
				return m, tea.Quit
			}

		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, min(msg.Height, m.maxHeight))
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *pickerModel[T]) View() string {
	if m.quitting {
		return ""
	}
	return m.list.View()
}
