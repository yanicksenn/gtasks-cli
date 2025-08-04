package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/gtasks/internal/version"
)

var rootCmd = &cobra.Command{
	Use:   "gtasks",
	Short: "A CLI for managing your Google Tasks",
	Long:  `gtasks is a powerful command-line interface that helps you manage your Google Tasks directly from the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		v, _ := cmd.Flags().GetBool("version")
		if v {
			fmt.Println(version.Get())
			os.Exit(0)
		}

		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().Bool("offline", false, "Enable offline mode")
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
