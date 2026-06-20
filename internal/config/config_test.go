package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/andrejmatveev/3xui-metrics-collector/internal/config"
)

func TestLoadValidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `
listen_addr: ":2112"
poll_interval: 30s
request_timeout: 10s
panels:
  - name: panel-a
    base_url: "https://example.com/panel"
    api_token: "secret"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.PollInterval != 30*time.Second {
		t.Fatalf("poll_interval: got %v", cfg.PollInterval)
	}
	if len(cfg.Panels) != 1 || cfg.Panels[0].Name != "panel-a" {
		t.Fatalf("panels: %+v", cfg.Panels)
	}
}

func TestValidateDuplicatePanelNames(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `
poll_interval: 30s
request_timeout: 10s
panels:
  - name: same
    base_url: "https://a.example.com"
    api_token: "t1"
  - name: same
    base_url: "https://b.example.com"
    api_token: "t2"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := config.Load(path); err == nil {
		t.Fatal("expected duplicate name error")
	}
}

func TestValidatePollInterval(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `
poll_interval: 1s
request_timeout: 10s
panels:
  - name: p
    base_url: "https://example.com"
    api_token: "t"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := config.Load(path); err == nil {
		t.Fatal("expected poll interval error")
	}
}

func TestEnvTokenOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `
poll_interval: 30s
request_timeout: 10s
panels:
  - name: my-panel
    base_url: "https://example.com"
    api_token: "from-file"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XUI_PANEL_MY_PANEL_TOKEN", "from-env")

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Panels[0].APIToken != "from-env" {
		t.Fatalf("token: got %q", cfg.Panels[0].APIToken)
	}
}
