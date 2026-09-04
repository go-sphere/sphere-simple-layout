package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewEmptyConfigProvidesUsableDefaults(t *testing.T) {
	config := NewEmptyConfig()
	if config.API.JWT == "" {
		t.Fatal("generated API JWT secret must be non-empty")
	}
	if got, want := config.API.HTTP.Address, "0.0.0.0:8899"; got != want {
		t.Errorf("API address = %q, want %q", got, want)
	}
}

func TestNewConfigValidatesAPISecret(t *testing.T) {
	_, err := NewConfig(writeConfig(t, `{}`))
	if err == nil {
		t.Fatal("NewConfig() error = nil, want missing JWT error")
	}
	if got, want := err.Error(), "api jwt must be non-empty"; got != want {
		t.Errorf("NewConfig() error = %q, want %q", got, want)
	}
}

func TestNewConfigAppliesLogLevelDefault(t *testing.T) {
	config, err := NewConfig(writeConfig(t, `{"api":{"jwt":"api-secret"}}`))
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}
	if got, want := config.Log.Level, "info"; got != want {
		t.Errorf("Log.Level = %q, want %q", got, want)
	}
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}
