package main

// An example demonstrating an application with multiple views.
//
// Note that this example was produced before the Bubbles progress component
// was available (github.com/charmbracelet/bubbles/progress) and thus, we're
// implementing a progress bar from scratch here.

import (
	"fmt"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fogleman/ease"
	"github.com/lucasb-eyer/go-colorful"
)

type (
	tickMsg  struct{}
	frameMsg struct{}
)

type appState struct {
	Choice   int
	Chosen   bool
	Ticks    int
	Frames   int
	Progress float64
	Loaded   bool
	Quitting bool
}

const (
	dotChar = " • "
)

// General stuff for styling the view
var (
	keywordStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	radioButtonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	dotStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle        = lipgloss.NewStyle().MarginLeft(2)
)

func main() {
	initialModel := appState{
		Choice:   0,
		Chosen:   false,
		Frames:   0,
		Progress: 0,
		Loaded:   false,
		Quitting: false,
	}
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

func (m appState) Init() tea.Cmd {
	return nil
}

// Main update function.
func (m appState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	// Hand off the message and appState to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.Chosen {
		return updateChoices(msg, m)
	}
	return updateChosen(msg, m)
}

// The main view, which just calls the appropriate sub-view
func (m appState) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	if !m.Chosen {
		s = choicesView(m)
	} else {
		s = chosenView(m)
	}
	return mainStyle.Render("\n" + s + "\n\n")
}

// Sub-update functions

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m appState) (tea.Model, tea.Cmd) {
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
			return m, frame()
		}
	}

	return m, nil
}

// Update loop for the second view after a choice has been made
func updateChosen(msg tea.Msg, m appState) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case frameMsg:
		if !m.Loaded {
			m.Frames++
			m.Progress = ease.OutBounce(float64(m.Frames) / float64(100))

			return m, frame()
		}
	}

	return m, nil
}

// The first view, where you're choosing a task
func choicesView(m appState) string {
	tpl := "nefs - nostr encrypted file sharing\n\n"
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

// The second view, after a task has been chosen
func chosenView(m appState) string {
	var msg string

	switch m.Choice {
	case 0:
		return sendView(m)
	case 1:
		return receiveView(m)
	default:
	}

	return msg + "\n\n" + "\n" + "%"
}

func sendView(m appState) string {
	return "send"
}

func receiveView(m appState) string {
	return "receive"
}

// Components
func radioButton(label string, checked bool) string {
	if checked {
		return radioButtonStyle.Render("(o) " + label)
	}
	return "( ) " + label
}

// Utils

// Convert a colorful.Color to a hexadecimal format.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}

