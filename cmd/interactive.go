package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start the interactive TUI",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running program: %w", err)
		}
		return nil
	},
}

type model struct {
	// TODO: Add fields for the TUI state
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return "Hello, World!"
}

func init() {
	RootCmd.AddCommand(interactiveCmd)

	// This is a temporary log file for debugging the TUI
	f, err := os.OpenFile("tui.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}
