package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
	"github.com/tgruben/slackit/internal/resolve"
)

var sendThread string

var sendCmd = &cobra.Command{
	Use:   "send <target> <message>",
	Short: "Send a message to a channel or user",
	Long: `Send a message to a Slack channel or user.

Target can be #channel-name, @username, @user@email.com, a raw Slack ID,
or a shortcut name from ~/.slackit.json.
Use "-" as the message to read from stdin.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := appConfig.ResolveShortcut(args[0])
		message := args[1]

		// Read from stdin if message is "-"
		if message == "-" {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
			message = strings.TrimRight(string(data), "\n")
		}

		channelID, err := resolve.Target(slackClient, target)
		if err != nil {
			return err
		}

		opts := []slack.MsgOption{
			slack.MsgOptionText(message, false),
		}
		if sendThread != "" {
			opts = append(opts, slack.MsgOptionTS(sendThread))
		}

		_, _, err = slackClient.PostMessage(channelID, opts...)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Message sent to %s\n", args[0])
		return nil
	},
}

func init() {
	sendCmd.Flags().StringVar(&sendThread, "thread", "", "Thread timestamp for threaded replies")
	rootCmd.AddCommand(sendCmd)
}
