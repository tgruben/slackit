package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/toddgruben/slackit/internal/format"
)

var (
	membersFormat      string
	membersIncludeBots bool
)

var membersCmd = &cobra.Command{
	Use:   "members",
	Short: "List workspace members",
	RunE: func(cmd *cobra.Command, args []string) error {
		outFmt, err := format.ParseFormat(membersFormat)
		if err != nil {
			return err
		}

		users, err := slackClient.GetUsers()
		if err != nil {
			return fmt.Errorf("failed to list users: %w", err)
		}

		type memberInfo struct {
			ID          string `json:"id"`
			Username    string `json:"username"`
			DisplayName string `json:"display_name"`
			RealName    string `json:"real_name"`
			Email       string `json:"email,omitempty"`
			IsBot       bool   `json:"is_bot,omitempty"`
		}

		var members []memberInfo
		for _, u := range users {
			if u.Deleted {
				continue
			}
			if u.IsBot && !membersIncludeBots {
				continue
			}
			// Skip Slackbot
			if u.ID == "USLACKBOT" && !membersIncludeBots {
				continue
			}
			members = append(members, memberInfo{
				ID:          u.ID,
				Username:    u.Name,
				DisplayName: u.Profile.DisplayName,
				RealName:    u.RealName,
				Email:       u.Profile.Email,
				IsBot:       u.IsBot,
			})
		}

		switch outFmt {
		case format.FormatTable:
			headers := []string{"ID", "USERNAME", "DISPLAY NAME", "REAL NAME", "EMAIL"}
			rows := make([][]string, len(members))
			for i, m := range members {
				rows[i] = []string{m.ID, m.Username, m.DisplayName, m.RealName, m.Email}
			}
			format.Table(os.Stdout, headers, rows)
		case format.FormatJSON:
			if err := format.JSON(os.Stdout, members); err != nil {
				return fmt.Errorf("failed to encode JSON: %w", err)
			}
		case format.FormatNames:
			names := make([]string, len(members))
			for i, m := range members {
				names[i] = m.Username
			}
			format.Names(os.Stdout, names)
		}

		return nil
	},
}

func init() {
	membersCmd.Flags().StringVar(&membersFormat, "format", "table", "Output format: table, json, or names")
	membersCmd.Flags().BoolVar(&membersIncludeBots, "include-bots", false, "Include bot users")
	rootCmd.AddCommand(membersCmd)
}
