package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Verify authentication and show workspace info",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := slackClient.AuthTest()
		if err != nil {
			return fmt.Errorf("auth test failed: %w", err)
		}

		fmt.Fprintf(os.Stdout, "Workspace:  %s\n", resp.Team)
		fmt.Fprintf(os.Stdout, "User:       %s\n", resp.User)
		fmt.Fprintf(os.Stdout, "User ID:    %s\n", resp.UserID)
		fmt.Fprintf(os.Stdout, "Team ID:    %s\n", resp.TeamID)
		fmt.Fprintf(os.Stdout, "URL:        %s\n", resp.URL)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
