package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/gtasks/internal/auth"
	"github.com/yanicksenn/gtasks/internal/config"
)

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage your authenticated accounts",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Google and add a new account",
	RunE: func(cmd *cobra.Command, args []string) error {
		user, err := auth.LoginViaWebFlow(context.Background())
		if err != nil {
			return fmt.Errorf("error during authentication: %w", err)
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}

		cfg.ActiveAccount = user
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}

		fmt.Printf("Successfully logged in as %s.\n", user)
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of the active account",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}

		if err := auth.Logout(cfg.ActiveAccount); err != nil {
			return fmt.Errorf("error logging out: %w", err)
		}

		fmt.Printf("Successfully logged out of %s.\n", cfg.ActiveAccount)
		cfg.ActiveAccount = ""
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}
		return nil
	},
}

var listAccountsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all authenticated accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		accounts, err := auth.ListAccounts()
		if err != nil {
			return fmt.Errorf("error listing accounts: %w", err)
		}

		if len(accounts) == 0 {
			fmt.Println("No accounts authenticated.")
			return nil
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}

		fmt.Println("Authenticated Accounts:")
		for _, account := range accounts {
			if account == cfg.ActiveAccount {
				fmt.Printf("- %s (active)\n", account)
			} else {
				fmt.Printf("- %s\n", account)
			}
		}
		return nil
	},
}

var switchAccountCmd = &cobra.Command{
	Use:   "switch [EMAIL]",
	Short: "Switch to a different account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		email := args[0]
		accounts, err := auth.ListAccounts()
		if err != nil {
			return fmt.Errorf("error listing accounts: %w", err)
		}

		for _, account := range accounts {
			if account == email {
				cfg, err := config.Load()
				if err != nil {
					return fmt.Errorf("error loading config: %w", err)
				}
				cfg.ActiveAccount = email
				if err := cfg.Save(); err != nil {
					return fmt.Errorf("error saving config: %w", err)
				}
				fmt.Printf("Successfully switched to %s.\n", email)
				return nil
			}
		}

		return fmt.Errorf("account %s not found. Please log in first", email)
	},
}

func init() {
	RootCmd.AddCommand(accountsCmd)
	accountsCmd.AddCommand(loginCmd)
	accountsCmd.AddCommand(logoutCmd)
	accountsCmd.AddCommand(listAccountsCmd)
	accountsCmd.AddCommand(switchAccountCmd)
}