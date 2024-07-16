package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	tickMsg  struct{}
	frameMsg struct{}
)

const (
	dotChar = " • "
)

type AppState struct {
	Choice   int
	Chosen   bool
	Ticks    int
	Frames   int
	Progress float64
	Loaded   bool
	Quitting bool

	SendState    SendState
	ReceiveState ReceiveState
}

func NewAppState() AppState {
	return AppState{
		Choice:       0,
		Chosen:       false,
		Frames:       0,
		Progress:     0,
		Loaded:       false,
		Quitting:     false,
		SendState:    NewSendState(),
		ReceiveState: NewReceiveState(),
	}
}

// General stuff for styling the view
var (
	keywordStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	radioButtonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	dotStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle        = lipgloss.NewStyle().MarginLeft(2)
)

func (m AppState) Init() tea.Cmd {
	return tea.Batch(
		m.SendState.Init(),
		m.ReceiveState.Init(),
	)
}

// Main update function.
func (m AppState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	if !m.Chosen {
		return updateChoices(msg, m)
	}

	var cmd tea.Cmd
	switch m.Choice {
	case 0:
		m.SendState, cmd = m.SendState.Update(msg)
	case 1:
		m.ReceiveState, cmd = m.ReceiveState.Update(msg)
	default:
	}

	return m, cmd
}

// The main view, which just calls the appropriate sub-view
func (m AppState) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}

	if !m.Chosen {
		s = choicesView(m)
	} else {
		switch m.Choice {
		case 0:
			s = m.SendState.View()
		case 1:
			s = m.ReceiveState.View()
		default:
		}
	}

	return mainStyle.Render("\n" + s + "\n\n")
}

// Sub-update functions

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m AppState) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice++
			if m.Choice > 1 {
				m.Choice = 1
			}
		case "k", "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = true
			return m, nil
		}
	}

	return m, nil
}

// The first view, where you're choosing a task
func choicesView(m AppState) string {
	tpl := "nefs - nostr encrypted file sharing\n"
	tpl += "%s\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")

	choices := fmt.Sprintf(
		"\n%s\n%s\n",
		radioButton("send a file", m.Choice == 0),
		radioButton("receive a file", m.Choice == 1),
	)

	return fmt.Sprintf(tpl, choices)
}
