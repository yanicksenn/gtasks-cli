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
		user, err := auth.LoginViaWebFlow(context.Background())
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

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of the active account",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		if err := auth.Logout(cfg.ActiveAccount); err != nil {
			fmt.Fprintf(os.Stderr, "Error logging out: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully logged out of %s.\n", cfg.ActiveAccount)
		cfg.ActiveAccount = ""
		if err := cfg.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
	},
}

var listAccountsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all authenticated accounts",
	Run: func(cmd *cobra.Command, args []string) {
		accounts, err := auth.ListAccounts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing accounts: %v\n", err)
			os.Exit(1)
		}

		if len(accounts) == 0 {
			fmt.Println("No accounts authenticated.")
			return
		}

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Authenticated Accounts:")
		for _, account := range accounts {
			if account == cfg.ActiveAccount {
				fmt.Printf("- %s (active)\n", account)
			} else {
				fmt.Printf("- %s\n", account)
			}
		}
	},
}

var switchAccountCmd = &cobra.Command{
	Use:   "switch [EMAIL]",
	Short: "Switch to a different account",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]
		accounts, err := auth.ListAccounts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing accounts: %v\n", err)
			os.Exit(1)
		}

		for _, account := range accounts {
			if account == email {
				cfg, err := config.Load()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
					os.Exit(1)
				}
				cfg.ActiveAccount = email
				if err := cfg.Save(); err != nil {
					fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("Successfully switched to %s.\n", email)
				return
			}
		}

		fmt.Fprintf(os.Stderr, "Account %s not found. Please log in first.\n", email)
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)
	accountsCmd.AddCommand(loginCmd)
	accountsCmd.AddCommand(logoutCmd)
	accountsCmd.AddCommand(listAccountsCmd)
	accountsCmd.AddCommand(switchAccountCmd)
}