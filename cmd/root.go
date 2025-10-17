package cmd

import (
	"fmt"
	"os"

	"github.com/onedusk/sb/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "sb",
	Short: "SB - Scalable media processing toolkit",
	Long: `SB is a powerful, extensible CLI tool for media processing.

Convert videos, images, and audio files with support for batch processing,
hardware acceleration, and quality controls.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// GetRootCmd returns the root command for adding subcommands
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sb.yaml)")
	rootCmd.PersistentFlags().IntP("workers", "w", 0, "number of parallel workers (default: CPU count)")
	rootCmd.PersistentFlags().StringP("out", "o", "", "output directory")
	rootCmd.PersistentFlags().BoolP("recursive", "r", false, "process directories recursively")
	rootCmd.PersistentFlags().BoolP("skip", "s", false, "skip existing files")
	rootCmd.PersistentFlags().BoolP("dry-run", "n", false, "preview without converting")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("flat", "f", false, "flatten output directory structure")

	// Bind flags to viper
	viper.BindPFlag("workers", rootCmd.PersistentFlags().Lookup("workers"))
	viper.BindPFlag("output_dir", rootCmd.PersistentFlags().Lookup("out"))
	viper.BindPFlag("skip_existing", rootCmd.PersistentFlags().Lookup("skip"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("flat_structure", rootCmd.PersistentFlags().Lookup("flat"))
}

// initConfig reads in config file and ENV variables
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	}

	// Initialize config
	if err := config.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}
}
