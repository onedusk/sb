package formats

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/onedusk/sb/internal/config"
	"github.com/onedusk/sb/internal/converter"
	"github.com/onedusk/sb/internal/processors/mov_to_mp4"
	"github.com/onedusk/sb/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// MP4 command flags
	mp4Quality      int
	mp4Preset       string
	mp4Codec        string
	mp4AudioCodec   string
	mp4AudioBitrate string
	mp4Bitrate      string
	mp4HWAccel      string
	mp4Dir          string
	mp4Recursive    bool
)

// MP4Cmd represents the mp4 command
var MP4Cmd = &cobra.Command{
	Use:   "mp4 [files...]",
	Short: "Convert video files to MP4 format",
	Long: `Convert video files to MP4 format using H.264/H.265 encoding.

Supports input formats: .mov, .avi, .mkv, .flv, .wmv, .m4v, .mpeg, .mpg, .webm

Examples:
  sb mp4 video.mov                       # Convert single file
  sb mp4 *.mov                           # Convert all MOV files
  sb mp4 -d ./videos -r                  # Convert directory recursively
  sb mp4 -q 20 -p slow video.mov         # High quality, slow preset
  sb mp4 --hw videotoolbox *.mov         # Hardware accelerated
  sb mp4 -w 8 -o ./converted *.mov       # 8 workers, custom output dir`,
	RunE: runMP4Convert,
}

func init() {
	// MP4-specific flags
	MP4Cmd.Flags().IntVarP(&mp4Quality, "quality", "q", 0, "CRF quality (0-51, lower = better, default: 23)")
	MP4Cmd.Flags().StringVarP(&mp4Preset, "preset", "p", "", "encoding preset (ultrafast|fast|medium|slow|veryslow)")
	MP4Cmd.Flags().StringVarP(&mp4Codec, "codec", "c", "", "video codec (h264|h265|vp9)")
	MP4Cmd.Flags().StringVar(&mp4AudioCodec, "audio", "", "audio codec (aac|mp3|copy)")
	MP4Cmd.Flags().StringVar(&mp4AudioBitrate, "audio-bitrate", "", "audio bitrate (e.g., 128k, 192k)")
	MP4Cmd.Flags().StringVarP(&mp4Bitrate, "bitrate", "b", "", "video bitrate (e.g., 2M, 5M)")
	MP4Cmd.Flags().StringVar(&mp4HWAccel, "hw", "", "hardware acceleration (videotoolbox|nvenc|qsv)")
	MP4Cmd.Flags().StringVarP(&mp4Dir, "dir", "d", "", "input directory")
	MP4Cmd.Flags().BoolVar(&mp4Recursive, "recursive", false, "process directory recursively")

	// Bind flags to viper with mp4 prefix
	viper.BindPFlag("mp4.quality", MP4Cmd.Flags().Lookup("quality"))
	viper.BindPFlag("mp4.preset", MP4Cmd.Flags().Lookup("preset"))
	viper.BindPFlag("mp4.codec", MP4Cmd.Flags().Lookup("codec"))
	viper.BindPFlag("mp4.audio", MP4Cmd.Flags().Lookup("audio"))
	viper.BindPFlag("mp4.bitrate", MP4Cmd.Flags().Lookup("bitrate"))
}

func runMP4Convert(cmd *cobra.Command, args []string) error {
	// Initialize config
	cfg := config.Get()

	// Gather input files
	inputs, err := gatherInputs(args, mp4Dir, mp4Recursive)
	if err != nil {
		return err
	}

	if len(inputs) == 0 {
		return fmt.Errorf("no input files found")
	}

	// Get converter
	conv, err := converter.Get("mp4")
	if err != nil {
		return fmt.Errorf("mp4 converter not available: %w", err)
	}

	mp4Conv, ok := conv.(*mov_to_mp4.MP4Converter)
	if !ok {
		return fmt.Errorf("invalid converter type")
	}

	// Build MP4 options
	mp4Opts := buildMP4Options(cfg)
	if err := mp4Conv.SetOptions(mp4Opts); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	// Build converter options
	convOpts := converter.Options{
		OutputDir:     viper.GetString("output_dir"),
		Workers:       viper.GetInt("workers"),
		SkipExisting:  viper.GetBool("skip_existing"),
		DryRun:        cmd.Flags().Changed("dry-run") && viper.GetBool("dry_run"),
		Verbose:       viper.GetBool("verbose"),
		FlatStructure: viper.GetBool("flat_structure"),
		ShowProgress:  true,
		Context:       context.Background(),
	}

	// Validate workers
	if convOpts.Workers <= 0 {
		convOpts.Workers = cfg.Workers
	}

	// Print conversion info
	if !convOpts.Verbose {
		ui.PrintInfo("Converting %d file(s) to MP4", len(inputs))
		ui.PrintInfo("Workers: %d, Quality: CRF %d, Preset: %s", convOpts.Workers, mp4Opts.CRF, mp4Opts.Preset)
		if mp4Opts.HWAccel != "" {
			ui.PrintInfo("Hardware acceleration: %s", mp4Opts.HWAccel)
		}
		fmt.Println()
	}

	// Convert files
	if len(inputs) == 1 {
		// Single file conversion
		result, err := mp4Conv.Convert(inputs[0], convOpts)
		if err != nil {
			return err
		}
		if !result.Skipped && !convOpts.DryRun {
			fmt.Printf("✓ %s -> %s\n", result.Input, result.Output)
		}
	} else {
		// Batch conversion
		_, err = mp4Conv.ConvertBatch(inputs, convOpts)
		if err != nil {
			return err
		}
	}

	return nil
}

// gatherInputs collects input files from args or directory
func gatherInputs(args []string, dir string, recursive bool) ([]string, error) {
	inputs := []string{}

	// If directory specified, scan it
	if dir != "" {
		if recursive {
			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && isVideoFile(path) {
					inputs = append(inputs, path)
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("error walking directory: %w", err)
			}
		} else {
			entries, err := os.ReadDir(dir)
			if err != nil {
				return nil, fmt.Errorf("error reading directory: %w", err)
			}
			for _, entry := range entries {
				if !entry.IsDir() {
					path := filepath.Join(dir, entry.Name())
					if isVideoFile(path) {
						inputs = append(inputs, path)
					}
				}
			}
		}
	}

	// Add files from args
	for _, arg := range args {
		// Check if it's a glob pattern
		matches, err := filepath.Glob(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", arg, err)
		}

		if len(matches) > 0 {
			for _, match := range matches {
				info, err := os.Stat(match)
				if err != nil {
					continue
				}
				if !info.IsDir() && isVideoFile(match) {
					inputs = append(inputs, match)
				}
			}
		} else {
			// Direct file path
			if info, err := os.Stat(arg); err == nil && !info.IsDir() {
				inputs = append(inputs, arg)
			}
		}
	}

	return inputs, nil
}

// isVideoFile checks if a file is a supported video format
func isVideoFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	videoExts := []string{".mov", ".avi", ".mkv", ".flv", ".wmv", ".m4v", ".mpeg", ".mpg", ".webm", ".mp4"}
	for _, validExt := range videoExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

// buildMP4Options builds MP4Options from config and flags
func buildMP4Options(cfg *config.Config) mov_to_mp4.MP4Options {
	opts := mov_to_mp4.DefaultMP4Options()

	// Apply config values
	if cfg.MP4.Quality > 0 {
		opts.CRF = cfg.MP4.Quality
	}
	if cfg.MP4.Preset != "" {
		opts.Preset = cfg.MP4.Preset
	}
	if cfg.MP4.Codec != "" {
		opts.VideoCodec = cfg.MP4.Codec
	}
	if cfg.MP4.Audio != "" {
		opts.AudioCodec = cfg.MP4.Audio
	}
	if cfg.MP4.Bitrate != "" {
		opts.VideoBitrate = cfg.MP4.Bitrate
	}
	if cfg.MP4.Hardware.Enabled && cfg.MP4.Hardware.Type != "" {
		opts.HWAccel = cfg.MP4.Hardware.Type
	}

	// Override with CLI flags
	if mp4Quality > 0 {
		opts.CRF = mp4Quality
	}
	if mp4Preset != "" {
		opts.Preset = mp4Preset
	}
	if mp4Codec != "" {
		opts.VideoCodec = mp4Codec
	}
	if mp4AudioCodec != "" {
		opts.AudioCodec = mp4AudioCodec
	}
	if mp4AudioBitrate != "" {
		opts.AudioBitrate = mp4AudioBitrate
	}
	if mp4Bitrate != "" {
		opts.VideoBitrate = mp4Bitrate
	}
	if mp4HWAccel != "" {
		opts.HWAccel = mp4HWAccel
	}

	return opts
}
