package auth

import (
	"fmt"
	"os"
)

// ResolveToken returns a Slack API token using priority:
// 1. Explicit token (from --token flag)
// 2. SLACK_TOKEN environment variable
// 3. ~/.slack/credentials.json (auto-selects if one workspace, requires workspace flag if multiple)
func ResolveToken(flagToken, workspace string) (string, error) {
	// Priority 1: explicit flag
	if flagToken != "" {
		return flagToken, nil
	}

	// Priority 2: environment variable
	if envToken := os.Getenv("SLACK_TOKEN"); envToken != "" {
		return envToken, nil
	}

	// Priority 3: credentials file
	creds, err := LoadCredentialsFile()
	if err != nil {
		return "", err
	}

	// If workspace specified, use it
	if workspace != "" {
		cred, ok := creds[workspace]
		if !ok {
			available := make([]string, 0, len(creds))
			for name := range creds {
				available = append(available, name)
			}
			return "", fmt.Errorf("workspace %q not found in credentials\nAvailable workspaces: %v", workspace, available)
		}
		return cred.Token, nil
	}

	// Auto-select if only one workspace
	if len(creds) == 1 {
		for _, cred := range creds {
			return cred.Token, nil
		}
	}

	// Multiple workspaces, no selection
	available := make([]string, 0, len(creds))
	for name := range creds {
		available = append(available, name)
	}
	return "", fmt.Errorf("multiple workspaces found, use --workspace to select one\nAvailable: %v", available)
}
