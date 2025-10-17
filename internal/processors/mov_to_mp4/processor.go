package mov_to_mp4

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/onedusk/sb/internal/converter"
	"github.com/onedusk/sb/internal/executor"
	"github.com/onedusk/sb/internal/ui"
)

func init() {
	// Auto-register this converter
	converter.Register(NewMP4Converter())
}

// MP4Converter implements the Converter interface for MOV to MP4 conversion
type MP4Converter struct {
	ffmpeg  *executor.FFmpeg
	options MP4Options
}

// NewMP4Converter creates a new MP4 converter
func NewMP4Converter() *MP4Converter {
	return &MP4Converter{
		options: DefaultMP4Options(),
	}
}

// Name returns the converter name
func (c *MP4Converter) Name() string {
	return "mp4"
}

// Description returns the converter description
func (c *MP4Converter) Description() string {
	return "Convert video files to MP4 format using H.264/H.265 encoding"
}

// SupportedInputs returns supported input formats
func (c *MP4Converter) SupportedInputs() []string {
	return []string{".mov", ".avi", ".mkv", ".flv", ".wmv", ".m4v", ".mpeg", ".mpg", ".webm"}
}

// OutputExtension returns the output extension
func (c *MP4Converter) OutputExtension() string {
	return ".mp4"
}

// Validate checks if the input file is valid
func (c *MP4Converter) Validate(input string) error {
	// Check if file exists
	info, err := os.Stat(input)
	if err != nil {
		return fmt.Errorf("cannot access file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("input is a directory, not a file")
	}

	// Check if extension is supported
	ext := strings.ToLower(filepath.Ext(input))
	supported := false
	for _, validExt := range c.SupportedInputs() {
		if ext == validExt {
			supported = true
			break
		}
	}

	if !supported {
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	return nil
}

// Convert processes a single file
func (c *MP4Converter) Convert(input string, opts converter.Options) (*converter.Result, error) {
	result := &converter.Result{
		Input: input,
	}

	start := time.Now()

	// Validate input
	if err := c.Validate(input); err != nil {
		result.Error = err
		result.Duration = time.Since(start)
		return result, err
	}

	// Initialize ffmpeg if needed
	if c.ffmpeg == nil {
		ff, err := executor.NewFFmpeg()
		if err != nil {
			result.Error = err
			result.Duration = time.Since(start)
			return result, err
		}
		c.ffmpeg = ff
	}

	// Determine output path
	output := c.determineOutputPath(input, opts)
	result.Output = output

	// Check if output already exists
	if opts.SkipExisting {
		if _, err := os.Stat(output); err == nil {
			result.Skipped = true
			result.SkipReason = "file already exists"
			result.Duration = time.Since(start)
			ui.PrintVerbose(opts.Verbose, "Skipping %s (already exists)", input)
			return result, nil
		}
	}

	// Get input file size
	if info, err := os.Stat(input); err == nil {
		result.InputSize = info.Size()
	}

	// Dry run mode
	if opts.DryRun {
		fmt.Printf("[DRY-RUN] Would convert: %s -> %s\n", input, output)
		result.Success = true
		result.Duration = time.Since(start)
		return result, nil
	}

	// Create output directory if needed
	outputDir := filepath.Dir(output)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		result.Error = fmt.Errorf("failed to create output directory: %w", err)
		result.Duration = time.Since(start)
		return result, result.Error
	}

	// Build ffmpeg options
	ffmpegOpts := executor.FFmpegOptions{
		VideoCodec:   c.options.VideoCodec,
		CRF:          c.options.CRF,
		Preset:       c.options.Preset,
		AudioCodec:   c.options.AudioCodec,
		AudioBitrate: c.options.AudioBitrate,
		HWAccel:      c.options.HWAccel,
		Bitrate:      c.options.VideoBitrate,
		Verbose:      opts.Verbose,
	}

	// Execute conversion
	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	ui.PrintVerbose(opts.Verbose, "Converting: %s -> %s", input, output)

	ffResult, err := c.ffmpeg.Convert(ctx, input, output, ffmpegOpts)
	result.Duration = time.Since(start)

	if err != nil {
		result.Error = fmt.Errorf("conversion failed: %w", err)
		ui.PrintVerbose(opts.Verbose, "Error: %v", err)
		if ffResult != nil && ffResult.Stderr != "" {
			ui.PrintVerbose(opts.Verbose, "FFmpeg stderr: %s", ffResult.Stderr)
		}
		return result, result.Error
	}

	// Get output file size
	if info, err := os.Stat(output); err == nil {
		result.OutputSize = info.Size()
	}

	result.Success = true
	ui.PrintVerbose(opts.Verbose, "Successfully converted %s in %s", input, result.Duration.Round(time.Millisecond))

	return result, nil
}

// ConvertBatch processes multiple files
func (c *MP4Converter) ConvertBatch(inputs []string, opts converter.Options) ([]*converter.Result, error) {
	if len(inputs) == 0 {
		return nil, fmt.Errorf("no input files provided")
	}

	results := make([]*converter.Result, 0, len(inputs))
	totalStart := time.Now()

	// Create progress bar
	showProgress := opts.ShowProgress && !opts.DryRun && !opts.Verbose
	progress := ui.NewProgressBar(len(inputs), "Converting", showProgress)

	// Create worker pool
	pool := executor.NewPool(opts.Workers)
	pool.Start()

	// Submit jobs
	go func() {
		for _, input := range inputs {
			inputCopy := input // Capture for closure
			pool.Submit(func(ctx context.Context) error {
				jobOpts := opts
				jobOpts.Context = ctx
				result, err := c.Convert(inputCopy, jobOpts)
				results = append(results, result)
				progress.Increment()

				if !opts.Verbose && !result.Skipped {
					if result.Success {
						fmt.Printf("✓ %s\n", inputCopy)
					} else {
						fmt.Printf("✗ %s: %v\n", inputCopy, err)
					}
				}

				return err
			})
		}
		pool.Stop()
	}()

	// Wait for completion
	errors := make([]error, 0)
	for err := range pool.Results() {
		if err != nil {
			errors = append(errors, err)
		}
	}

	progress.Finish()

	// Calculate statistics
	totalDuration := time.Since(totalStart)
	success := 0
	failed := 0
	skipped := 0

	for _, result := range results {
		if result.Skipped {
			skipped++
		} else if result.Success {
			success++
		} else {
			failed++
		}
	}

	if !opts.Verbose {
		ui.PrintSummary(len(inputs), success, failed, skipped, totalDuration)
	}

	if len(errors) > 0 {
		return results, fmt.Errorf("%d conversion(s) failed", len(errors))
	}

	return results, nil
}

// SetOptions sets converter-specific options
func (c *MP4Converter) SetOptions(opts MP4Options) error {
	if err := opts.Validate(); err != nil {
		return err
	}
	c.options = opts
	return nil
}

// determineOutputPath calculates the output file path
func (c *MP4Converter) determineOutputPath(input string, opts converter.Options) string {
	baseName := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
	outputName := baseName + c.OutputExtension()

	// If output directory is specified
	if opts.OutputDir != "" {
		if opts.FlatStructure {
			// Flat structure: output all files to output dir
			return filepath.Join(opts.OutputDir, outputName)
		} else {
			// Preserve structure: maintain relative path
			// For now, just use the output dir with filename
			// TODO: implement full path preservation
			return filepath.Join(opts.OutputDir, outputName)
		}
	}

	// No output dir specified: place next to input file
	return filepath.Join(filepath.Dir(input), outputName)
}
