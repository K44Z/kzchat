package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
