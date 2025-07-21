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
	token := &oauth2.Token{AccessToken: "test-token"}
	err := cache.save("test@example.com", token)
	if err != nil {
		t.Fatalf("failed to save token cache: %v", err)
	}

	// Verify the file was created.
	expectedPath := filepath.Join(tempDir, ".config", tokenFile)
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Fatalf("token cache file was not created at %s", expectedPath)
	}

	// Test loading the token cache.
	loadedCache, err := loadTokenCache()
	if err != nil {
		t.Fatalf("failed to load token cache: %v", err)
	}

	loadedToken, ok := loadedCache.Tokens["test@example.com"]
	if !ok {
		t.Fatalf("token for 'test@example.com' not found in cache")
	}

	if loadedToken.AccessToken != "test-token" {
		t.Errorf("expected access token to be 'test-token', got '%s'", loadedToken.AccessToken)
	}
}
