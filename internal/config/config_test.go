package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	// Create a temporary directory for the test.
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// Test saving a config.
	cfg := &Config{ActiveAccount: "test@example.com"}
	err := cfg.Save()
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Verify the file was created.
	expectedPath := filepath.Join(tempDir, ".config", "gtasks.yml")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Fatalf("config file was not created at %s", expectedPath)
	}

	// Test loading the config.
	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loadedCfg.ActiveAccount != "test@example.com" {
		t.Errorf("expected active account to be 'test@example.com', got '%s'", loadedCfg.ActiveAccount)
	}
}
