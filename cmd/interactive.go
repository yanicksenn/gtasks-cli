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
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		tasklist, _ := cmd.Flags().GetString("tasklist")
		if tasklist == "" {
			return fmt.Errorf("the --tasklist flag is required")
		}

		m, err := tui.New(h.Client, tasklist)
		if err != nil {
			return fmt.Errorf("error creating new model: %w", err)
		}

		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running program: %w", err)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(interactiveCmd)
	interactiveCmd.Flags().String("tasklist", "@default", "The ID of the task list")

	// This is a temporary log file for debugging the TUI
	f, err := os.OpenFile("tui.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}