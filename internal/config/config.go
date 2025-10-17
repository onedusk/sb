package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	// Global settings
	Workers       int    `mapstructure:"workers"`
	SkipExisting  bool   `mapstructure:"skip_existing"`
	OutputDir     string `mapstructure:"output_dir"`
	FlatStructure bool   `mapstructure:"flat_structure"`
	Verbose       bool   `mapstructure:"verbose"`

	// Format-specific settings
	MP4 MP4Config `mapstructure:"mp4"`
}

// MP4Config contains MP4-specific configuration
type MP4Config struct {
	Quality   int            `mapstructure:"quality"` // CRF value
	Preset    string         `mapstructure:"preset"`
	Codec     string         `mapstructure:"codec"`
	Audio     string         `mapstructure:"audio"`
	Bitrate   string         `mapstructure:"bitrate"`
	Hardware  HardwareConfig `mapstructure:"hardware"`
}

// HardwareConfig contains hardware acceleration settings
type HardwareConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Type    string `mapstructure:"type"` // videotoolbox, nvenc, qsv
}

var (
	defaultConfig *Config
)

// Initialize sets up Viper and loads configuration
func Initialize() error {
	// Set config name and paths
	viper.SetConfigName(".sb")
	viper.SetConfigType("yaml")

	// Add config paths
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(home)
	}
	viper.AddConfigPath(".")

	// Set defaults
	setDefaults()

	// Enable environment variables
	viper.SetEnvPrefix("SB")
	viper.AutomaticEnv()

	// Try to read config file (not an error if it doesn't exist)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal into config struct
	defaultConfig = &Config{}
	if err := viper.Unmarshal(defaultConfig); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	return nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Global defaults
	viper.SetDefault("workers", runtime.NumCPU())
	viper.SetDefault("skip_existing", false)
	viper.SetDefault("output_dir", "")
	viper.SetDefault("flat_structure", false)
	viper.SetDefault("verbose", false)

	// MP4 defaults
	viper.SetDefault("mp4.quality", 23)
	viper.SetDefault("mp4.preset", "medium")
	viper.SetDefault("mp4.codec", "h264")
	viper.SetDefault("mp4.audio", "aac")
	viper.SetDefault("mp4.bitrate", "")
	viper.SetDefault("mp4.hardware.enabled", false)
	viper.SetDefault("mp4.hardware.type", "")
}

// Get returns the current configuration
func Get() *Config {
	if defaultConfig == nil {
		Initialize()
	}
	return defaultConfig
}

// GetInt gets an integer config value
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetString gets a string config value
func GetString(key string) string {
	return viper.GetString(key)
}

// GetBool gets a boolean config value
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// Set sets a config value
func Set(key string, value interface{}) {
	viper.Set(key, value)
}

// WriteConfig writes the current configuration to file
func WriteConfig(path string) error {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = filepath.Join(home, ".sb.yaml")
	}

	return viper.WriteConfigAs(path)
}

// ExampleConfig returns an example configuration file content
func ExampleConfig() string {
	return `# SB Media Processor Configuration

# Global settings
workers: 4              # Number of parallel workers (default: CPU count)
skip_existing: false    # Skip files that already exist
output_dir: ""          # Default output directory (empty = same as input)
flat_structure: false   # Flatten directory structure in output
verbose: false          # Enable verbose logging

# MP4 conversion settings
mp4:
  quality: 23           # CRF value (0-51, lower = better quality)
  preset: medium        # Encoding preset (ultrafast, fast, medium, slow, veryslow)
  codec: h264           # Video codec (h264, h265, vp9)
  audio: aac            # Audio codec (aac, mp3, copy)
  bitrate: ""           # Video bitrate (e.g., "2M", "5M", empty = use CRF)
  hardware:
    enabled: false      # Enable hardware acceleration
    type: ""            # Hardware type (videotoolbox, nvenc, qsv)

# Future format settings can be added here
# jpg:
#   quality: 95
# png:
#   compression: 9
`
}
