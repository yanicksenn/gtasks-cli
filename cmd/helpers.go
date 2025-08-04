package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/gtasks/internal/config"
	"github.com/yanicksenn/gtasks/internal/gtasks"
	"github.com/yanicksenn/gtasks/internal/ui"
)

// CommandHelper encapsulates common command dependencies.
type CommandHelper struct {
	Client  gtasks.Client
	Config  *config.Config
	Printer *ui.Printer
}

// NewCommandHelper creates a new CommandHelper.
func NewCommandHelper(cmd *cobra.Command) (*CommandHelper, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	client, err := gtasks.NewClientFromCommand(cmd, context.Background())
	if err != nil {
		return nil, err
	}

	quiet, _ := cmd.Flags().GetBool("quiet")
	outputFormat, _ := cmd.Flags().GetString("output")

	printer := ui.NewPrinter(cmd.OutOrStdout(), outputFormat, quiet)

	return &CommandHelper{
		Client:  client,
		Config:  cfg,
		Printer: printer,
	}, nil
}

