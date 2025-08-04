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
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}
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

		return h.Printer.PrintSuccess(fmt.Sprintf("Successfully logged in as %s.", user))
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of the active account",
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}

		if err := auth.Logout(cfg.ActiveAccount); err != nil {
			return fmt.Errorf("error logging out: %w", err)
		}

		h.Printer.PrintSuccess(fmt.Sprintf("Successfully logged out of %s.", cfg.ActiveAccount))
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
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		accounts, err := auth.ListAccounts()
		if err != nil {
			return fmt.Errorf("error listing accounts: %w", err)
		}

		return h.Printer.PrintAccounts(accounts, h.Config.ActiveAccount)
	},
}

var switchAccountCmd = &cobra.Command{
	Use:   "switch [EMAIL]",
	Short: "Switch to a different account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}
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
				return h.Printer.PrintSuccess(fmt.Sprintf("Successfully switched to %s.", email))
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