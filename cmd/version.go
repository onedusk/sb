package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print version, commit, and build date information for sb.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sb version %s\n", Version)
		fmt.Printf("commit: %s\n", Commit)
		fmt.Printf("built: %s\n", Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
