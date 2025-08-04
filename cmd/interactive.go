package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/yanicksenn/gtasks/internal/tui"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start the interactive TUI",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(tui.InitialModel())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running program: %w", err)
		}
		return nil
	},
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