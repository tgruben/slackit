package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
	"github.com/tgruben/slackit/internal/resolve"
)

var (
	uploadComment  string
	uploadFilename string
)

var uploadCmd = &cobra.Command{
	Use:   "upload <target> <filepath>",
	Short: "Upload a file to a channel or user",
	Long: `Upload a file to a Slack channel or user.

Target can be #channel-name, @username, @user@email.com, a raw Slack ID,
or a shortcut name from ~/.slackit.json.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := appConfig.ResolveShortcut(args[0])
		filePath := args[1]

		// Verify file exists
		info, err := os.Stat(filePath)
		if err != nil {
			return fmt.Errorf("cannot access file %q: %w", filePath, err)
		}
		if info.IsDir() {
			return fmt.Errorf("%q is a directory, not a file", filePath)
		}

		channelID, err := resolve.Target(slackClient, target)
		if err != nil {
			return err
		}

		filename := filepath.Base(filePath)
		if uploadFilename != "" {
			filename = uploadFilename
		}

		params := slack.UploadFileParameters{
			Channel:        channelID,
			File:           filePath,
			Filename:       filename,
			Title:          filename,
			InitialComment: uploadComment,
		}

		_, err = slackClient.UploadFile(params)
		if err != nil {
			return fmt.Errorf("failed to upload file: %w", err)
		}

		fmt.Fprintf(os.Stderr, "File %s uploaded to %s\n", filename, args[0])
		return nil
	},
}

func init() {
	uploadCmd.Flags().StringVar(&uploadComment, "comment", "", "Initial comment for the uploaded file")
	uploadCmd.Flags().StringVar(&uploadFilename, "filename", "", "Override the filename")
	rootCmd.AddCommand(uploadCmd)
}
