package cmd

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
	"github.com/tgruben/slackit/internal/auth"
	"github.com/tgruben/slackit/internal/config"
)

var (
	flagToken     string
	flagWorkspace string
	flagDebug     bool

	slackClient *slack.Client
	appConfig   *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "slackit",
	Short: "A CLI tool for posting messages and files to Slack",
	Long:  "Slackit lets you send messages, upload files, and manage Slack workspaces from the command line.",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load config (non-fatal if missing)
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		appConfig = cfg

		// Skip auth for commands that don't need it
		if cmd.Name() == "version" || cmd.Name() == "help" {
			return nil
		}

		token, err := auth.ResolveToken(flagToken, flagWorkspace)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		opts := []slack.Option{}
		if flagDebug {
			opts = append(opts, slack.OptionDebug(true))
		}

		slackClient = slack.New(token, opts...)

		// Verify the token works
		resp, err := slackClient.AuthTest()
		if err != nil {
			return fmt.Errorf("auth test failed: %w\nCheck your token or run `slack login`", err)
		}

		if flagDebug {
			fmt.Printf("Authenticated as %s in %s\n", resp.User, resp.Team)
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagToken, "token", "t", "", "Slack API token (overrides env/credentials)")
	rootCmd.PersistentFlags().StringVarP(&flagWorkspace, "workspace", "w", "", "Workspace name from ~/.slack/credentials.json")
	rootCmd.PersistentFlags().BoolVar(&flagDebug, "debug", false, "Enable debug output")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
