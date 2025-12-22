package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	argc := len(os.Args)
	if argc != 2 {
		log.Fatal("Please provide your save's filepath to the executable")
	}

	saveFilePath := os.Args[1]
	p := tea.NewProgram(initialModel(saveFilePath))

	_, err := p.Run()
	handle_err(err)
}
