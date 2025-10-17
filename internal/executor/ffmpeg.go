package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// FFmpegOptions contains options for ffmpeg execution
type FFmpegOptions struct {
	// Video options
	VideoCodec    string // h264, h265, vp9
	CRF           int    // Constant Rate Factor (0-51, lower = better quality)
	Preset        string // ultrafast, superfast, veryfast, faster, fast, medium, slow, slower, veryslow
	Bitrate       string // e.g., "2M", "5M"

	// Audio options
	AudioCodec    string // aac, mp3, copy
	AudioBitrate  string // e.g., "128k", "192k"

	// Hardware acceleration
	HWAccel       string // videotoolbox, nvenc, qsv
	HWAccelDevice string // optional device specification

	// Advanced
	ExtraArgs     []string
	Verbose       bool
}

// FFmpegResult contains the result of an ffmpeg execution
type FFmpegResult struct {
	Success  bool
	Duration time.Duration
	Stdout   string
	Stderr   string
	Error    error
}

// FFmpeg wraps ffmpeg command execution
type FFmpeg struct {
	binaryPath string
}

// NewFFmpeg creates a new FFmpeg executor
func NewFFmpeg() (*FFmpeg, error) {
	// Check if ffmpeg is available
	path, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}

	return &FFmpeg{
		binaryPath: path,
	}, nil
}

// Convert executes ffmpeg to convert a file
func (f *FFmpeg) Convert(ctx context.Context, input, output string, opts FFmpegOptions) (*FFmpegResult, error) {
	args := f.buildArgs(input, output, opts)

	if opts.Verbose {
		fmt.Printf("[ffmpeg] %s %s\n", f.binaryPath, strings.Join(args, " "))
	}

	start := time.Now()

	cmd := exec.CommandContext(ctx, f.binaryPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(start)

	result := &FFmpegResult{
		Success:  err == nil,
		Duration: duration,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Error:    err,
	}

	return result, err
}

// GetInfo retrieves media file information using ffprobe
func (f *FFmpeg) GetInfo(ctx context.Context, input string) (string, error) {
	probePath := strings.Replace(f.binaryPath, "ffmpeg", "ffprobe", 1)

	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		input,
	}

	cmd := exec.CommandContext(ctx, probePath, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffprobe failed: %w", err)
	}

	return string(output), nil
}

// buildArgs constructs ffmpeg command arguments
func (f *FFmpeg) buildArgs(input, output string, opts FFmpegOptions) []string {
	args := []string{"-y"} // Always overwrite output files

	// Hardware acceleration (must come before input)
	if opts.HWAccel != "" {
		args = append(args, "-hwaccel", opts.HWAccel)
		if opts.HWAccelDevice != "" {
			args = append(args, "-hwaccel_device", opts.HWAccelDevice)
		}
	}

	// Input file
	args = append(args, "-i", input)

	// Video codec
	if opts.VideoCodec != "" {
		codec := opts.VideoCodec
		// Map to hardware-accelerated codec if requested
		if opts.HWAccel == "videotoolbox" {
			switch codec {
			case "h264":
				codec = "h264_videotoolbox"
			case "h265", "hevc":
				codec = "hevc_videotoolbox"
			}
		}
		args = append(args, "-c:v", codec)
	} else {
		args = append(args, "-c:v", "libx264")
	}

	// Quality settings
	if opts.CRF > 0 {
		// CRF only works with certain codecs
		if !strings.Contains(opts.VideoCodec, "videotoolbox") {
			args = append(args, "-crf", fmt.Sprintf("%d", opts.CRF))
		} else {
			// VideoToolbox uses different quality scale
			args = append(args, "-q:v", fmt.Sprintf("%d", opts.CRF))
		}
	}

	// Preset (encoding speed vs compression)
	if opts.Preset != "" && !strings.Contains(opts.VideoCodec, "videotoolbox") {
		args = append(args, "-preset", opts.Preset)
	}

	// Bitrate (overrides CRF if both specified)
	if opts.Bitrate != "" {
		args = append(args, "-b:v", opts.Bitrate)
	}

	// Audio codec
	if opts.AudioCodec != "" {
		args = append(args, "-c:a", opts.AudioCodec)
	} else {
		args = append(args, "-c:a", "aac")
	}

	// Audio bitrate
	if opts.AudioBitrate != "" {
		args = append(args, "-b:a", opts.AudioBitrate)
	}

	// Extra arguments
	if len(opts.ExtraArgs) > 0 {
		args = append(args, opts.ExtraArgs...)
	}

	// Output file
	args = append(args, output)

	return args
}

// CheckVersion returns the ffmpeg version
func (f *FFmpeg) CheckVersion() (string, error) {
	cmd := exec.Command(f.binaryPath, "-version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return lines[0], nil
	}

	return "", fmt.Errorf("unable to parse ffmpeg version")
}
