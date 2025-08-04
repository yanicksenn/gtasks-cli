package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yanicksenn/gtasks/internal/version"
)

var RootCmd = &cobra.Command{
	Use:   "gtasks",
	Short: "A CLI for managing your Google Tasks",
	Long:  `gtasks is a powerful command-line interface that helps you manage your Google Tasks directly from the terminal.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		v, _ := cmd.Flags().GetBool("version")
		if v {
			cmd.Println(version.Get())
			return nil
		}

		if len(args) == 0 {
			return cmd.Help()
		}
		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().Bool("offline", false, "Enable offline mode")
	RootCmd.PersistentFlags().BoolP("quiet", "q", false, "Disable output")
	RootCmd.PersistentFlags().StringP("output", "o", "table", "Output format (table, json, yaml)")
	RootCmd.Flags().BoolP("version", "v", false, "Print the version number")
}

func Execute() error {
	return RootCmd.Execute()
}
