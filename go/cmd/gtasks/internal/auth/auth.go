package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	tokenFile = "gtasks-token.json"
)

// ErrCredentialsNotFound is returned when the user's credentials are not found.
var ErrCredentialsNotFound = errors.New("credentials not found. Please run 'gtasks login'")

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

	conf := &oauth2.Config{
		ClientID:     "1021942592516-ddskqoqs4d752kpmrak83vmsq05k5n07.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/tasks", "https://www.googleapis.com/auth/tasks.readonly", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

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

// TokenCache represents the structure of the credentials file.
type TokenCache struct {
	Tokens map[string]*oauth2.Token `json:"tokens"`
}

func getTokenCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", tokenFile), nil
}

func loadTokenCache() (*TokenCache, error) {
	path, err := getTokenCachePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &TokenCache{Tokens: make(map[string]*oauth2.Token)}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cache TokenCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	return &cache, nil
}

func (tc *TokenCache) save(user string, token *oauth2.Token) error {
	tc.Tokens[user] = token
	return tc.saveAll()
}

func (tc *TokenCache) saveAll() error {
	path, err := getTokenCachePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(tc, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
