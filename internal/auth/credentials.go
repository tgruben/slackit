package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Credential represents a single workspace credential from the Slack CLI.
type Credential struct {
	Token string `json:"token"`
	Team  string `json:"team"`
}

// CredentialsFile maps workspace names to credentials.
type CredentialsFile map[string]Credential

// LoadCredentialsFile reads and parses ~/.slack/credentials.json.
func LoadCredentialsFile() (CredentialsFile, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot determine home directory: %w", err)
	}

	path := filepath.Join(home, ".slack", "credentials.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w\nInstall the Slack CLI or use --token / SLACK_TOKEN", path, err)
	}

	var creds CredentialsFile
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("cannot parse %s: %w", path, err)
	}

	if len(creds) == 0 {
		return nil, fmt.Errorf("no workspaces found in %s\nRun `slack login` to authenticate", path)
	}

	return creds, nil
}
