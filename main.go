package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

var Program *tea.Program

func main() {
	m := NewModel()
	Program = tea.NewProgram(&m, tea.WithAltScreen())
	if err := Program.Start(); err != nil {
		log.Fatal(err)
	}
}
