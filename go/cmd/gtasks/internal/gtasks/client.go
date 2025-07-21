package gtasks

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/workspace/go/cmd/gtasks/internal/auth"
	"github.com/yanicksenn/workspace/go/cmd/gtasks/internal/config"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

// Client is a wrapper for the Google Tasks API.
type Client struct {
	service TasksService
}

// NewClient creates a new client for the Google Tasks API.
// It will create an online client by default, unless the --offline flag is set.
func NewClient(cmd *cobra.Command, ctx context.Context) (*Client, error) {
	offline, _ := cmd.Flags().GetBool("offline")
	if offline {
		// TODO: Return an offline client
		return nil, nil
	}

	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	authenticator, err := auth.NewAuthenticator()
	if err != nil {
		return nil, err
	}

	httpClient, err := authenticator.GetClient(ctx, cfg.ActiveAccount)
	if err != nil {
		return nil, err
	}

	service, err := tasks.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	return &Client{service: &TasksServiceWrapper{service: service}}, nil
}