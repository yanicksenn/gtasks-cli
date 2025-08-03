package auth

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/oauth2"
)

func TestTokenCache(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// Test saving a token.
	cache := &TokenCache{Tokens: make(map[string]*oauth2.Token)}
	token1 := &oauth2.Token{AccessToken: "test-token-1"}
	token2 := &oauth2.Token{AccessToken: "test-token-2"}
	err := cache.save("test1@example.com", token1)
	if err != nil {
		t.Fatalf("failed to save token cache: %v", err)
	}
	err = cache.save("test2@example.com", token2)
	if err != nil {
		t.Fatalf("failed to save token cache: %v", err)
	}

	// Verify the file was created.
	expectedPath := filepath.Join(tempDir, ".config", TokenFile)
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Fatalf("token cache file was not created at %s", expectedPath)
	}

	// Test listing accounts.
	accounts, err := ListAccounts()
	if err != nil {
		t.Fatalf("failed to list accounts: %v", err)
	}
	if len(accounts) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(accounts))
	}

	// Test logging out.
	err = Logout("test1@example.com")
	if err != nil {
		t.Fatalf("failed to logout: %v", err)
	}

	// Verify the account was removed.
	accounts, err = ListAccounts()
	if err != nil {
		t.Fatalf("failed to list accounts: %v", err)
	}
	if len(accounts) != 1 {
		t.Fatalf("expected 1 account, got %d", len(accounts))
	}
	if accounts[0] == "test1@example.com" {
		t.Fatalf("account was not removed")
	}
}
