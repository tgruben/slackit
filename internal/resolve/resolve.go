package resolve

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
)

// Target resolves a user-provided target string to a Slack channel ID.
//
// Supported formats:
//   - #channel-name — looks up channel by name
//   - @username — looks up user by display name, then opens a DM
//   - @user@email.com — looks up user by email, then opens a DM
//   - C.../D.../G... — used directly as channel ID
//   - U... — opens a DM with the user ID
func Target(client *slack.Client, target string) (string, error) {
	switch {
	case strings.HasPrefix(target, "#"):
		return resolveChannel(client, strings.TrimPrefix(target, "#"))
	case strings.HasPrefix(target, "@"):
		return resolveUser(client, strings.TrimPrefix(target, "@"))
	case strings.HasPrefix(target, "C"), strings.HasPrefix(target, "D"), strings.HasPrefix(target, "G"):
		return target, nil
	case strings.HasPrefix(target, "U"):
		return openDM(client, target)
	default:
		return "", fmt.Errorf("unrecognized target %q\nUse #channel, @user, @email, or a Slack ID (C.../U...)", target)
	}
}

func resolveChannel(client *slack.Client, name string) (string, error) {
	var cursor string
	for {
		params := &slack.GetConversationsParameters{
			Types:           []string{"public_channel", "private_channel"},
			Limit:           200,
			Cursor:          cursor,
			ExcludeArchived: true,
		}
		channels, nextCursor, err := client.GetConversations(params)
		if err != nil {
			return "", fmt.Errorf("failed to list channels: %w", err)
		}
		for _, ch := range channels {
			if ch.Name == name {
				return ch.ID, nil
			}
		}
		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}
	return "", fmt.Errorf("channel #%s not found\nUse `slackit channels` to list available channels", name)
}

func resolveUser(client *slack.Client, identifier string) (string, error) {
	// Check if it's an email
	if strings.Contains(identifier, "@") {
		user, err := client.GetUserByEmail(identifier)
		if err != nil {
			return "", fmt.Errorf("user with email %q not found: %w", identifier, err)
		}
		return openDM(client, user.ID)
	}

	// Search by display name / real name
	users, err := client.GetUsers()
	if err != nil {
		return "", fmt.Errorf("failed to list users: %w", err)
	}
	for _, u := range users {
		if u.Deleted {
			continue
		}
		if strings.EqualFold(u.Name, identifier) ||
			strings.EqualFold(u.Profile.DisplayName, identifier) ||
			strings.EqualFold(u.RealName, identifier) {
			return openDM(client, u.ID)
		}
	}
	return "", fmt.Errorf("user @%s not found\nUse `slackit members` to list workspace members", identifier)
}

func openDM(client *slack.Client, userID string) (string, error) {
	params := &slack.OpenConversationParameters{
		Users: []string{userID},
	}
	channel, _, _, err := client.OpenConversation(params)
	if err != nil {
		return "", fmt.Errorf("failed to open DM with %s: %w", userID, err)
	}
	return channel.ID, nil
}
