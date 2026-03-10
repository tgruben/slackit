package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
	"github.com/toddgruben/slackit/internal/format"
)

var (
	channelsFormat          string
	channelsIncludeArchived bool
)

var channelsCmd = &cobra.Command{
	Use:   "channels",
	Short: "List workspace channels",
	RunE: func(cmd *cobra.Command, args []string) error {
		outFmt, err := format.ParseFormat(channelsFormat)
		if err != nil {
			return err
		}

		type channelInfo struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Topic    string `json:"topic,omitempty"`
			Members  int    `json:"members"`
			Archived bool   `json:"archived,omitempty"`
		}

		var channels []channelInfo
		var cursor string
		for {
			params := &slack.GetConversationsParameters{
				Types:           []string{"public_channel", "private_channel"},
				Limit:           200,
				Cursor:          cursor,
				ExcludeArchived: !channelsIncludeArchived,
			}
			convs, nextCursor, err := slackClient.GetConversations(params)
			if err != nil {
				return fmt.Errorf("failed to list channels: %w", err)
			}
			for _, ch := range convs {
				channels = append(channels, channelInfo{
					ID:       ch.ID,
					Name:     ch.Name,
					Topic:    ch.Topic.Value,
					Members:  ch.NumMembers,
					Archived: ch.IsArchived,
				})
			}
			if nextCursor == "" {
				break
			}
			cursor = nextCursor
		}

		switch outFmt {
		case format.FormatTable:
			headers := []string{"ID", "NAME", "MEMBERS", "TOPIC"}
			rows := make([][]string, len(channels))
			for i, ch := range channels {
				topic := ch.Topic
				if len(topic) > 50 {
					topic = topic[:47] + "..."
				}
				rows[i] = []string{ch.ID, ch.Name, strconv.Itoa(ch.Members), topic}
			}
			format.Table(os.Stdout, headers, rows)
		case format.FormatJSON:
			if err := format.JSON(os.Stdout, channels); err != nil {
				return fmt.Errorf("failed to encode JSON: %w", err)
			}
		case format.FormatNames:
			names := make([]string, len(channels))
			for i, ch := range channels {
				names[i] = ch.Name
			}
			format.Names(os.Stdout, names)
		}

		return nil
	},
}

func init() {
	channelsCmd.Flags().StringVar(&channelsFormat, "format", "table", "Output format: table, json, or names")
	channelsCmd.Flags().BoolVar(&channelsIncludeArchived, "include-archived", false, "Include archived channels")
	rootCmd.AddCommand(channelsCmd)
}
