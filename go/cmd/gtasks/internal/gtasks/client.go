package gtasks

import (
	"context"

	"github.com/yanicksenn/workspace/go/cmd/gtasks/internal/auth"
	"github.com/yanicksenn/workspace/go/cmd/gtasks/internal/config"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

// Client is a wrapper for the Google Tasks API.
type Client struct {
	service TasksService
}

// NewClient creates a new authenticated client for the Google Tasks API.
func NewClient(ctx context.Context) (*Client, error) {
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
