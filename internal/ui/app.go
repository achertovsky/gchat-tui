package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var conversations = []string{
	"Engineering",
	"Design Team",
	"Ada Lovelace",
	"Project Updates",
}

type model struct {
	cursor   int
	selected int
	width    int
	height   int
}

// Run starts the terminal UI.
func Run() error {
	_, err := tea.NewProgram(model{selected: -1}, tea.WithAltScreen()).Run()
	return err
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(conversations)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.cursor
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	sidebarWidth := min(28, max(20, m.width/3))
	contentWidth := max(1, m.width-sidebarWidth)
	bodyHeight := max(1, m.height-2)

	sidebar := m.sidebar(sidebarWidth, bodyHeight)
	content := m.content(contentWidth, bodyHeight)
	footer := lipgloss.NewStyle().
		Width(m.width).
		Foreground(lipgloss.Color("241")).
		Render("↑/k up  ↓/j down  enter select  q quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content),
		footer,
	)
}

func (m model) sidebar(width, height int) string {
	var items []string
	for index, conversation := range conversations {
		prefix := "  "
		if index == m.cursor {
			prefix = "> "
		}
		items = append(items, prefix+conversation)
	}

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder(), false, true, false, false).
		BorderForeground(lipgloss.Color("240")).
		Padding(1).
		Render("Conversations\n\n" + strings.Join(items, "\n"))
}

func (m model) content(width, height int) string {
	body := "Choose a conversation and press Enter."
	title := "Messages"
	if m.selected >= 0 {
		title = conversations[m.selected]
		body = fmt.Sprintf(
			"Placeholder messages for %s\n\nNo messages have been loaded yet.",
			conversations[m.selected],
		)
	}

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2).
		Render(title + "\n\n" + body)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
