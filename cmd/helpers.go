package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/gtasks/internal/config"
	"github.com/yanicksenn/gtasks/internal/gtasks"
)

// CommandHelper encapsulates common command dependencies.
type CommandHelper struct {
	Client gtasks.Client
	Config *config.Config
	Quiet  bool
}

// NewCommandHelper creates a new CommandHelper.
func NewCommandHelper(cmd *cobra.Command) (*CommandHelper, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	client, err := gtasks.NewClient(cmd, context.Background())
	if err != nil {
		return nil, err
	}

	quiet, _ := cmd.Flags().GetBool("quiet")

	return &CommandHelper{
		Client: client,
		Config: cfg,
		Quiet:  quiet,
	}, nil
}
