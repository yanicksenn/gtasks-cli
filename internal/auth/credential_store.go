package auth

import (
	"encoding/json"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

// TokenFile is the name of the file where the OAuth2 tokens are stored.
const TokenFile = "gtasks-token.json"

// TokenCache represents the structure of the credentials file.
type TokenCache struct {
	Tokens map[string]*oauth2.Token `json:"tokens"`
}

// getTokenCachePath returns the path to the token cache file.
func getTokenCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", TokenFile), nil
}

// loadTokenCache loads the token cache from the file system.
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

// save saves the token for the given user to the cache.
func (tc *TokenCache) save(user string, token *oauth2.Token) error {
	tc.Tokens[user] = token
	return tc.saveAll()
}

// saveAll saves the entire token cache to the file system.
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
