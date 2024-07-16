package view

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type SendState struct {
	filepicker   filepicker.Model
	selectedFile string
	quitting     bool
	err          error
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func NewSendState() SendState {
	fp := filepicker.New()
	fp.CurrentDirectory, _ = os.UserHomeDir()

	return SendState{
		filepicker: fp,
	}
}

func (s SendState) Init() tea.Cmd {
	return s.filepicker.Init()
}

func (s SendState) Update(msg tea.Msg) (SendState, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			s.quitting = true
			return s, tea.Quit
		}
	case clearErrorMsg:
		s.err = nil
	default:
	}

	var cmd tea.Cmd
	s.filepicker, cmd = s.filepicker.Update(msg)

	if didSelect, path := s.filepicker.DidSelectFile(msg); didSelect {
		s.selectedFile = path
	}

	if didSelect, path := s.filepicker.DidSelectDisabledFile(msg); didSelect {
		s.err = errors.New(path + " is not valid")
		s.selectedFile = ""

		return s, tea.Batch(cmd, clearErrorAfter(3*time.Second))
	}

	return s, cmd
}

func (s SendState) View() string {
	if s.quitting {
		return "bye!"
	}

	var view strings.Builder
	view.WriteString("\n ")

	if s.err != nil {
		view.WriteString(s.filepicker.Styles.DisabledFile.Render("oops: " + s.err.Error()))
	} else if s.selectedFile == "" {
		view.WriteString("pick a file: ")
	} else {
		view.WriteString("selected file: " + s.filepicker.Styles.Selected.Render(s.selectedFile))
	}

	view.WriteString("\n\n" + s.filepicker.View() + "\n")

	return view.String()
}
