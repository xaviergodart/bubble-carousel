package carousel

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model defines a state for the carousel widget.
type Model struct {
	KeyMap KeyMap

	items        []string
	cursor       int
	width        int
	height       int
	focus        bool
	evenlySpaced bool
	styles       Styles

	content   string
	itemWidth int
	start     int
	end       int
}

// KeyMap defines keybindings. It satisfies to the help.KeyMap interface, which
// is used to render the menu.
type KeyMap struct {
	SelectLeft  key.Binding
	SelectRight key.Binding
}

// DefaultKeyMap returns a default set of keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		SelectLeft: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "h"),
		),
		SelectRight: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
	}
}

// Styles contains style definitions for this carousel component. By default,
// vthese alues are generated by DefaultStyles.
type Styles struct {
	Item     lipgloss.Style
	Selected lipgloss.Style
}

// DefaultStyles returns a set of default style definitions for this carousel.
func DefaultStyles() Styles {
	return Styles{
		Item: lipgloss.NewStyle().Padding(0, 1),
		Selected: lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("212")),
	}
}

// SetStyles sets the table styles.
func (m *Model) SetStyles(s Styles) {
	m.styles = s
	m.UpdateSize()
}

// Option is used to set options in New. For example:
//
//	carousel := New(WithItems([]string{"Item 1", "Item 2", "Item 3"}))
type Option func(*Model)

// New creates a new model for the carousel widget.
func New(opts ...Option) Model {
	m := Model{
		cursor: 0,

		KeyMap: DefaultKeyMap(),
		styles: DefaultStyles(),
	}

	for _, opt := range opts {
		opt(&m)
	}

	m.UpdateSize()

	return m
}

// WithItems sets the carousel items (data).
func WithItems(items []string) Option {
	return func(m *Model) {
		m.SetItems(items)
	}
}

// WithEvenlySpacedItems sets all items with the same width.
func WithEvenlySpacedItems() Option {
	return func(m *Model) {
		m.evenlySpaced = true
	}
}

// WithHeight sets the height of the carousel.
func WithHeight(h int) Option {
	return func(m *Model) {
		m.height = h
	}
}

// WithWidth sets the width of the carousel.
func WithWidth(w int) Option {
	return func(m *Model) {
		m.width = w
	}
}

// WithFocused sets the focus state of the carousel.
func WithFocused(f bool) Option {
	return func(m *Model) {
		m.focus = f
	}
}

// WithStyles sets the carousel styles.
func WithStyles(s Styles) Option {
	return func(m *Model) {
		m.styles = s
	}
}

// WithKeyMap sets the key map.
func WithKeyMap(km KeyMap) Option {
	return func(m *Model) {
		m.KeyMap = km
	}
}

// Update is the Bubble Tea update loop.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focus {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.SelectLeft):
			m.MoveLeft()
		case key.Matches(msg, m.KeyMap.SelectRight):
			m.MoveRight()
		}
	}

	return m, nil
}

// Focused returns the focus state of the carousel.
func (m Model) Focused() bool {
	return m.focus
}

// Focus focuses the carousel, allowing the user to move around the items and
// interact.
func (m *Model) Focus() {
	m.focus = true
	m.UpdateSize()
}

// Blur blurs the carousel, preventing selection or movement.
func (m *Model) Blur() {
	m.focus = false
	m.UpdateSize()
}

// View renders the component.
func (m Model) View() string {
	return m.content
}

// UpdateSize updates the carousel size based on the previously defined
// items and width.
func (m *Model) UpdateSize() {
	items := make([]string, 0, len(m.items))
	width := 0
	m.end = len(m.items)

	for i := range m.items {
		item := m.renderItem(i)

		if i >= m.start {
			width += lipgloss.Width(item)
		}

		items = append(items, item)

		if i == m.cursor && width > m.width {
			m.start++
		}

		if i == m.cursor && i < m.start {
			m.start--

			width += lipgloss.Width(item)
		}

		if i > m.cursor && width > m.width {
			m.end = i

			break
		}
	}

	m.content = lipgloss.JoinHorizontal(
		lipgloss.Center,
		items[m.start:m.end]...,
	)
}

// SelectedItem returns the selected item.
func (m Model) SelectedItem() string {
	return m.items[m.cursor]
}

// Items returns the current items.
func (m Model) Items() []string {
	return m.items
}

// SetItems sets a new items state.
func (m *Model) SetItems(items []string) {
	m.items = items
	m.itemWidth = 0
	m.cursor = 0
	m.start = 0
	for i := range m.items {
		item := m.renderItem(i)
		m.itemWidth = max(m.itemWidth, lipgloss.Width(item))
	}
	m.UpdateSize()
}

// SetWidth sets the width of the carousel.
func (m *Model) SetWidth(w int) {
	m.width = w
	m.UpdateSize()
}

// SetHeight sets the height of the carousel.
func (m *Model) SetHeight(h int) {
	m.height = h
	m.UpdateSize()
}

// Height returns the height of the carousel.
func (m Model) Height() int {
	return m.height
}

// Width returns the width of the carousel.
func (m Model) Width() int {
	return m.width
}

// Cursor returns the index of the selected row.
func (m Model) Cursor() int {
	return m.cursor
}

// HasRightItems returns true if there's items left on the right.
func (m Model) HasRightItems() bool {
	return m.end < len(m.items)
}

// HasLeftItems returns true if there's items left on the left.
func (m Model) HasLeftItems() bool {
	return m.start > 0
}

// SetCursor sets the cursor position in the carousel.
func (m *Model) SetCursor(n int) {
	m.cursor = clamp(n, 0, len(m.items)-1)
	m.UpdateSize()
}

// MoveLeft moves the selection left by one item..
// It can not go before the first item.
func (m *Model) MoveLeft() {
	m.cursor = clamp(m.cursor-1, 0, len(m.items)-1)
	m.UpdateSize()
}

// MoveDown moves the selection right by one item.
// It can not go after the last row.
func (m *Model) MoveRight() {
	m.cursor = clamp(m.cursor+1, 0, len(m.items)-1)
	m.UpdateSize()
}

func (m *Model) renderItem(itemID int) string {
	var item string
	if itemID == m.cursor {
		item = m.styles.Selected.Render(m.items[itemID])
	} else {
		item = m.styles.Item.Render(m.items[itemID])
	}

	if m.evenlySpaced {
		return lipgloss.Place(
			m.itemWidth,
			m.height,
			lipgloss.Left,
			lipgloss.Center,
			item,
		)
	}

	return item
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}
