package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const minPollInterval = 5 * time.Second

type Config struct {
	ListenAddr     string        `yaml:"listen_addr"`
	PollInterval   time.Duration `yaml:"poll_interval"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
	Panels         []PanelConfig `yaml:"panels"`
}

type PanelConfig struct {
	Name               string `yaml:"name"`
	BaseURL            string `yaml:"base_url"`
	APIToken           string `yaml:"api_token"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
	CollectOutbounds   bool   `yaml:"collect_outbounds"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	cfg := &Config{
		ListenAddr:     ":2112",
		PollInterval:   30 * time.Second,
		RequestTimeout: 10 * time.Second,
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	applyEnvOverrides(cfg)
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	for i := range cfg.Panels {
		envKey := "XUI_PANEL_" + strings.ToUpper(sanitizeEnvName(cfg.Panels[i].Name)) + "_TOKEN"
		if token := os.Getenv(envKey); token != "" {
			cfg.Panels[i].APIToken = token
		}
	}
}

func sanitizeEnvName(name string) string {
	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	return b.String()
}

func (c *Config) Validate() error {
	if c.ListenAddr == "" {
		return fmt.Errorf("listen_addr must not be empty")
	}
	if c.PollInterval < minPollInterval {
		return fmt.Errorf("poll_interval must be at least %s", minPollInterval)
	}
	if c.RequestTimeout <= 0 {
		return fmt.Errorf("request_timeout must be positive")
	}
	if len(c.Panels) == 0 {
		return fmt.Errorf("at least one panel is required")
	}

	seen := make(map[string]struct{}, len(c.Panels))
	for i, p := range c.Panels {
		if p.Name == "" {
			return fmt.Errorf("panels[%d]: name must not be empty", i)
		}
		if p.BaseURL == "" {
			return fmt.Errorf("panels[%d] %q: base_url must not be empty", i, p.Name)
		}
		if _, err := url.Parse(p.BaseURL); err != nil {
			return fmt.Errorf("panels[%d] %q: invalid base_url: %w", i, p.Name, err)
		}
		if p.APIToken == "" {
			return fmt.Errorf("panels[%d] %q: api_token must not be empty", i, p.Name)
		}
		if _, ok := seen[p.Name]; ok {
			return fmt.Errorf("duplicate panel name %q", p.Name)
		}
		seen[p.Name] = struct{}{}
	}
	return nil
}
