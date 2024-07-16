package view

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ReceiveState struct {
	Err error
}

func NewReceiveState() ReceiveState {
	return ReceiveState{}
}

func (s ReceiveState) Init() tea.Cmd {
	return nil
}

func (s ReceiveState) Update(msg tea.Msg) (ReceiveState, tea.Cmd) {
	return s, nil
}

func (s ReceiveState) View() string {
	return "receive"
}
