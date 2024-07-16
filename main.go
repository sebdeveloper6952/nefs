package main

// An example demonstrating an application with multiple views.
//
// Note that this example was produced before the Bubbles progress component
// was available (github.com/charmbracelet/bubbles/progress) and thus, we're
// implementing a progress bar from scratch here.

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sebdeveloper6952/nefs/view"
)

func main() {
	p := tea.NewProgram(view.NewAppState())
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
