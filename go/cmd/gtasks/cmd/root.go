package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gtasks",
	Short: "A CLI for managing your Google Tasks",
	Long:  `gtasks is a powerful command-line interface that helps you manage your Google Tasks directly from the terminal.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
