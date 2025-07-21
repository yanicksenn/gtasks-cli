package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/workspace/go/cmd/gtasks/internal/auth"
	"github.com/yanicksenn/workspace/go/cmd/gtasks/internal/config"
)

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage your authenticated accounts",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Google and add a new account",
	Run: func(cmd *cobra.Command, args []string) {
		authenticator, err := auth.NewAuthenticator()
		if err != nil {
			// This error is likely due to missing credentials.json
			fmt.Println("To authenticate with Google Tasks, you need to create your own OAuth 2.0 Client ID.")
			fmt.Println("Please visit the following URL to create your credentials:")
			fmt.Println("https://developers.google.com/tasks/docs/create-credentials")
			fmt.Println("\nOnce you have your credentials, please create a file at `~/.config/credentials.json` with the following format:")
			fmt.Println(`{
  "installed": {
    "client_id": "YOUR_CLIENT_ID",
    "client_secret": "YOUR_CLIENT_SECRET"
  }
}`)
			fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
			os.Exit(1)
		}

		_, user, err := authenticator.NewClient(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during authentication: %v\n", err)
			os.Exit(1)
		}

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		cfg.ActiveAccount = user
		if err := cfg.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully logged in as %s.\n", user)
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)
	accountsCmd.AddCommand(loginCmd)
}
