package auth

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// ErrCredentialsNotFound is returned when the user's credentials are not found.
var ErrCredentialsNotFound = errors.New("credentials not found. Please run 'gtasks login'")

// getOAuthConfig returns the OAuth2 configuration for the application.
func getOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "1021942592516-ddskqoqs4d752kpmrak83vmsq05k5n07.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/tasks", "https://www.googleapis.com/auth/tasks.readonly", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

// GetClient returns an authenticated HTTP client for the given user.
func GetClient(ctx context.Context, user string) (*http.Client, error) {
	cache, err := loadTokenCache()
	if err != nil {
		return nil, err
	}

	token, ok := cache.Tokens[user]
	if !ok {
		return nil, ErrCredentialsNotFound
	}

	conf := getOAuthConfig()
	return oauth2.NewClient(ctx, conf.TokenSource(ctx, token)), nil
}

// Logout removes the token for the given user from the cache.
func Logout(user string) error {
	cache, err := loadTokenCache()
	if err != nil {
		return err
	}

	delete(cache.Tokens, user)

	return cache.saveAll()
}

// ListAccounts lists all the accounts in the token cache.
func ListAccounts() ([]string, error) {
	cache, err := loadTokenCache()
	if err != nil {
		return nil, err
	}

	var accounts []string
	for user := range cache.Tokens {
		accounts = append(accounts, user)
	}

	return accounts, nil
}