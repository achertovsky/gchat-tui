package ui

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

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

// Message contains presentation-ready data for a conversation message.
type Message struct {
	SenderName string
	Text       string
	Timestamp  time.Time
}

// MessagePage contains a chronological page and a token for older messages.
type MessagePage struct {
	Messages      []Message
	NextPageToken string
}

// MessageLoader retrieves a page of messages for a conversation.
type MessageLoader func(context.Context, string, string) (MessagePage, error)

type model struct {
	loader         ConversationLoader
	messageLoader  MessageLoader
	conversations  []Conversation
	messages       []Message
	cursor         int
	selected       int
	loading        bool
	loadError      error
	messageLoading bool
	messageError   error
	nextPageToken  string
	messageRequest int
	scrollOffset   int
	width          int
	height         int
}

// Run starts the terminal UI.
func Run(loader ConversationLoader, messageLoader MessageLoader) error {
	_, err := tea.NewProgram(model{
		loader:        loader,
		messageLoader: messageLoader,
		selected:      -1,
		loading:       true,
	}, tea.WithAltScreen()).Run()
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
				m.messages = nil
				m.nextPageToken = ""
				m.messageError = nil
				m.messageLoading = true
				m.scrollOffset = 0
				m.messageRequest++
				return m, m.loadMessages("", false)
			}
		case "r":
			m.loading = true
			m.loadError = nil
			return m, m.loadConversations()
		case "u":
			if m.selected >= 0 && m.nextPageToken != "" && !m.messageLoading {
				m.messageLoading = true
				m.messageError = nil
				m.messageRequest++
				return m, m.loadMessages(m.nextPageToken, true)
			}
		case "pgup", "ctrl+u":
			m.scrollOffset = max(0, m.scrollOffset-5)
		case "pgdown", "ctrl+d":
			m.scrollOffset += 5
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
	case messagesLoadedMsg:
		if msg.request != m.messageRequest {
			return m, nil
		}
		m.messageLoading = false
		m.messageError = msg.err
		if msg.err == nil {
			if msg.older {
				m.messages = append(msg.page.Messages, m.messages...)
			} else {
				m.messages = msg.page.Messages
			}
			sort.SliceStable(m.messages, func(left, right int) bool {
				if m.messages[left].Timestamp.IsZero() {
					return false
				}
				if m.messages[right].Timestamp.IsZero() {
					return true
				}
				return m.messages[left].Timestamp.Before(m.messages[right].Timestamp)
			})
			m.nextPageToken = msg.page.NextPageToken
			if msg.older {
				m.scrollOffset = 0
			}
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
		Render("↑/k up  ↓/j down  enter open  u older  pgup/pgdn scroll  r refresh  q quit")

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
	title := "Messages"
	if m.selected < 0 || m.selected >= len(m.conversations) {
		return m.renderContent(width, height, title, []string{"Choose a conversation and press Enter."})
	}

	title = displayName(m.conversations[m.selected])
	if m.messageLoading {
		return m.renderContent(width, height, title, []string{"Loading messages..."})
	}
	if m.messageError != nil {
		return m.renderContent(width, height, title, []string{"Unable to load messages.", "", "Press Enter to retry."})
	}
	if len(m.messages) == 0 {
		return m.renderContent(width, height, title, []string{"No messages in this conversation."})
	}

	lines := renderMessages(m.messages)
	if m.nextPageToken != "" {
		lines = append(lines, "", "Press u to load older messages.")
	}
	return m.renderContent(width, height, title, lines)
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

type messagesLoadedMsg struct {
	page    MessagePage
	older   bool
	request int
	err     error
}

func (m model) loadMessages(pageToken string, older bool) tea.Cmd {
	conversationName := m.conversations[m.selected].Name
	requestID := m.messageRequest
	return func() tea.Msg {
		if m.messageLoader == nil {
			return messagesLoadedMsg{request: requestID, err: fmt.Errorf("message loading is not configured")}
		}
		page, err := m.messageLoader(context.Background(), conversationName, pageToken)
		return messagesLoadedMsg{page: page, older: older, request: requestID, err: err}
	}
}

func (m model) renderContent(width, height int, title string, lines []string) string {
	availableLines := max(1, height-4)
	maxOffset := max(0, len(lines)-availableLines)
	offset := min(m.scrollOffset, maxOffset)
	end := min(len(lines), offset+availableLines)
	body := strings.Join(lines[offset:end], "\n")
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2).
		Render(title + "\n\n" + body)
}

func renderMessages(messages []Message) []string {
	lines := make([]string, 0, len(messages)*3)
	for _, message := range messages {
		sender := strings.TrimSpace(message.SenderName)
		if sender == "" {
			sender = "Unknown sender"
		}
		timestamp := "Unknown time"
		if !message.Timestamp.IsZero() {
			timestamp = message.Timestamp.Local().Format("2006-01-02 15:04")
		}
		text := strings.TrimSpace(message.Text)
		if text == "" {
			text = "[Message unavailable]"
		}
		lines = append(lines, fmt.Sprintf("%s  %s", sender, timestamp), text, "")
	}
	return lines
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
