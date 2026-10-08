package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewEmptyConfigProvidesUsableDefaults(t *testing.T) {
	config := NewEmptyConfig()
	if got, want := config.API.HTTP.Address, "0.0.0.0:8899"; got != want {
		t.Errorf("API address = %q, want %q", got, want)
	}
}

func TestNewConfigAppliesLogLevelDefault(t *testing.T) {
	config, err := NewConfig(writeConfig(t, `{}`))
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

func TestNewConfigLoadsHTTPServerOptions(t *testing.T) {
	config, err := NewConfig(writeConfig(t, `{"api":{"http":{
		"trusted_proxies":["10.0.0.0/8"],"max_body_bytes":1024,"read_timeout_seconds":5,"idle_timeout_seconds":-1}}}`))
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}
	got := config.API.HTTP.Options
	if len(got.TrustedProxies) != 1 || got.TrustedProxies[0] != "10.0.0.0/8" ||
		got.MaxBodyBytes != 1024 || got.ReadTimeoutSeconds != 5 || got.IdleTimeoutSeconds != -1 {
		t.Errorf("API.HTTP options = %+v", got)
	}
}

func TestNewConfigRejectsInvalidTrustedProxy(t *testing.T) {
	_, err := NewConfig(writeConfig(t, `{"api":{"http":{"trusted_proxies":["proxy.internal"]}}}`))
	if err == nil || !strings.Contains(err.Error(), "api http: trusted_proxies") {
		t.Fatalf("NewConfig() error = %v, want a api http trusted_proxies error", err)
	}
}
