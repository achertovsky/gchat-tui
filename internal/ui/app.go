package ui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Conversation is a sidebar item supplied by the application.
type Conversation struct {
	Name        string
	DisplayName string
	Type        string
}

// ConversationLoader retrieves the conversations displayed by the TUI.
type ConversationLoader func(context.Context) ([]Conversation, error)

type model struct {
	loader        ConversationLoader
	conversations []Conversation
	cursor        int
	selected      int
	loading       bool
	loadError     error
	width         int
	height        int
}

// Run starts the terminal UI.
func Run(loader ConversationLoader) error {
	_, err := tea.NewProgram(model{loader: loader, selected: -1, loading: true}, tea.WithAltScreen()).Run()
	return err
}

func (m model) Init() tea.Cmd {
	return m.loadConversations()
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
			if m.cursor < len(m.conversations)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.conversations) > 0 {
				m.selected = m.cursor
			}
		case "r":
			m.loading = true
			m.loadError = nil
			return m, m.loadConversations()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case conversationsLoadedMsg:
		m.loading = false
		m.loadError = msg.err
		if msg.err == nil {
			m.conversations = msg.conversations
			m.cursor = 0
			m.selected = -1
		}
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
		Render("↑/k up  ↓/j down  enter select  r refresh  q quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content),
		footer,
	)
}

func (m model) sidebar(width, height int) string {
	var items []string
	for index, conversation := range m.conversations {
		prefix := "  "
		if index == m.cursor {
			prefix = "> "
		}
		items = append(items, prefix+conversationLabel(conversation))
	}
	body := "Loading conversations..."
	if m.loadError != nil {
		body = "Unable to load conversations.\n\nPress r to retry."
	} else if !m.loading && len(items) == 0 {
		body = "No conversations found."
	} else if len(items) > 0 {
		body = strings.Join(items, "\n")
	}

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder(), false, true, false, false).
		BorderForeground(lipgloss.Color("240")).
		Padding(1).
		Render("Conversations\n\n" + body)
}

func (m model) content(width, height int) string {
	body := "Choose a conversation and press Enter."
	title := "Messages"
	if m.selected >= 0 && m.selected < len(m.conversations) {
		title = displayName(m.conversations[m.selected])
		body = fmt.Sprintf(
			"Placeholder messages for %s\n\nNo messages have been loaded yet.",
			title,
		)
	}

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2).
		Render(title + "\n\n" + body)
}

type conversationsLoadedMsg struct {
	conversations []Conversation
	err           error
}

func (m model) loadConversations() tea.Cmd {
	return func() tea.Msg {
		if m.loader == nil {
			return conversationsLoadedMsg{err: fmt.Errorf("conversation loading is not configured")}
		}
		conversations, err := m.loader(context.Background())
		return conversationsLoadedMsg{conversations: conversations, err: err}
	}
}

func displayName(conversation Conversation) string {
	if strings.TrimSpace(conversation.DisplayName) != "" {
		return conversation.DisplayName
	}
	if strings.TrimSpace(conversation.Name) != "" {
		return conversation.Name
	}
	return "Unnamed conversation"
}

func conversationLabel(conversation Conversation) string {
	label := "[?]"
	switch conversation.Type {
	case "DIRECT_MESSAGE":
		label = "[DM]"
	case "GROUP_CHAT":
		label = "[GC]"
	case "SPACE":
		label = "[S]"
	}
	return label + " " + displayName(conversation)
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
