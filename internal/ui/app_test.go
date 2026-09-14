package ui

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func TestOpeningConversationFocusesComposer(t *testing.T) {
	input := textinput.New()
	initial := model{
		conversations: []Conversation{{Name: "spaces/example"}},
		input:         input,
	}

	updated, _ := initial.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := updated.(model)
	if !got.input.Focused() {
		t.Fatal("opening a conversation should focus the composer")
	}
}

func TestComposeKeyFocusesComposer(t *testing.T) {
	input := textinput.New()
	initial := model{
		conversations: []Conversation{{Name: "spaces/example"}},
		selected:      0,
		input:         input,
	}

	updated, _ := initial.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("i")})
	got := updated.(model)
	if !got.input.Focused() {
		t.Fatal("compose key should focus the composer")
	}
}
