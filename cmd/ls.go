package cmd

import (
	"fmt"
	"strings"

	"github.com/onedusk/sb/internal/converter"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List available converters",
	Long:  `List all available converters and their supported formats.`,
	Run: func(cmd *cobra.Command, args []string) {
		converters := converter.ListConverters()

		if len(converters) == 0 {
			fmt.Println("No converters available")
			return
		}

		fmt.Println("Available Converters:")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()

		for _, conv := range converters {
			fmt.Printf("  %s\n", conv.Name())
			fmt.Printf("    Description: %s\n", conv.Description())
			fmt.Printf("    Input:       %s\n", strings.Join(conv.SupportedInputs(), ", "))
			fmt.Printf("    Output:      %s\n", conv.OutputExtension())
			fmt.Println()
		}

		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("Total: %d converter(s)\n", len(converters))
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
