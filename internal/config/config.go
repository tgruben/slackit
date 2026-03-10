package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the slackit configuration file.
type Config struct {
	Channels map[string]string `json:"channels"`
}

// Load reads ~/.slackit.json, returning an empty config if the file doesn't exist.
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return &Config{}, nil
	}

	path := filepath.Join(home, ".slackit.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse %s: %w", path, err)
	}

	return &cfg, nil
}

// ResolveShortcut checks if the target matches a configured channel shortcut.
// Returns the channel ID if found, or the original target if not.
func (c *Config) ResolveShortcut(target string) string {
	if c == nil || c.Channels == nil {
		return target
	}
	if id, ok := c.Channels[target]; ok {
		return id
	}
	return target
}
