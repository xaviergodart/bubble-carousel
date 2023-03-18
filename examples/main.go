package main

import (
	"fmt"
	"os"

	carousel "github.com/xaviergodart/bubble-carousel"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	Margin(1, 1)

type model struct {
	carousel carousel.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.carousel.SetWidth(msg.Width - 2)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.carousel.Focused() {
				m.carousel.Blur()
			} else {
				m.carousel.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.carousel.SelectedItem()),
			)
		}
	}
	m.carousel, cmd = m.carousel.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var left, right string
	if m.carousel.HasLeftItems() {
		left = "◀"
	}
	if m.carousel.HasRightItems() {
		right = "▶"
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		left,
		baseStyle.Render(m.carousel.View()),
		right,
	)
}

func main() {
	nb := 20
	items := make([]string, 0, nb)
	for i := 0; i < nb; i++ {
		items = append(items, fmt.Sprintf("ITEM %d", i+1))
	}

	t := carousel.New(
		carousel.WithItems(items),
		carousel.WithFocused(true),
		carousel.WithEvenlySpacedItems(),
	)

	s := carousel.DefaultStyles()
	s.Item = s.Item.Padding(1, 1)
	s.Selected = s.Selected.
		Padding(1, 1).
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57"))
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
