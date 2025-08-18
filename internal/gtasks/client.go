package gtasks

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/gtasks/internal/auth"
	"github.com/yanicksenn/gtasks/internal/config"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

// onlineClient is a client that interacts with the real Google Tasks API.
// It implements the Client interface.
type onlineClient struct {
	service *tasks.Service
}

// NewClientFromCommand creates a new client based on the --offline flag from a cobra command.
func NewClientFromCommand(cmd *cobra.Command, ctx context.Context) (Client, error) {
	offline, _ := cmd.Flags().GetBool("offline")
	return NewClient(ctx, offline)
}

// NewClient creates a new client based on the offline flag.
func NewClient(ctx context.Context, offline bool) (Client, error) {
	if offline {
		return newOfflineClient()
	}
	return newOnlineClient(ctx)
}

func newOnlineClient(ctx context.Context) (*onlineClient, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	var httpClient *http.Client
	if cfg.ActiveAccount != "" {
		var getClientErr error
		httpClient, getClientErr = auth.GetClient(ctx, cfg.ActiveAccount)
		if getClientErr != nil {
			if errors.Is(getClientErr, auth.ErrTokenRefreshFailed) || errors.Is(getClientErr, auth.ErrCredentialsNotFound) {
				httpClient = nil
			} else {
				return nil, getClientErr
			}
		}
	}

	if httpClient == nil {
		fmt.Println("Authentication required. Please follow the instructions to log in.")
		user, loginErr := auth.LoginViaWebFlow(ctx)
		if loginErr != nil {
			return nil, fmt.Errorf("authentication failed: %w", loginErr)
		}

		cfg.ActiveAccount = user
		if err := cfg.Save(); err != nil {
			return nil, fmt.Errorf("failed to save active account: %w", err)
		}

		var getClientErr error
		httpClient, getClientErr = auth.GetClient(ctx, user)
		if getClientErr != nil {
			return nil, fmt.Errorf("failed to get client after login: %w", getClientErr)
		}
	}

	service, err := tasks.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	return &onlineClient{service: service}, nil
}
