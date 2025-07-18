package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/K44Z/kzchat/configs"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

var Program *tea.Program

func main() {
	envPath := filepath.Join(string(os.PathSeparator), "etc", "kzchat", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("error loading env file from :", envPath)
	}
	configs.SetBaseUrl()
	m := NewModel()
	Program = tea.NewProgram(&m, tea.WithAltScreen())
	if err := Program.Start(); err != nil {
		log.Fatal(err)
	}
}
