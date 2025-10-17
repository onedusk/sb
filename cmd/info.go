package cmd

import (
	"context"
	"fmt"

	"github.com/onedusk/sb/internal/executor"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info [file]",
	Short: "Display media file information",
	Long:  `Display detailed information about a media file using ffprobe.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]

		ffmpeg, err := executor.NewFFmpeg()
		if err != nil {
			return fmt.Errorf("ffmpeg not available: %w", err)
		}

		fmt.Printf("Media Info: %s\n", input)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		info, err := ffmpeg.GetInfo(context.Background(), input)
		if err != nil {
			return fmt.Errorf("failed to get info: %w", err)
		}

		fmt.Println(info)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
