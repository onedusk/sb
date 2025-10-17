package ui

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

// ProgressBar wraps progressbar with additional features
type ProgressBar struct {
	bar     *progressbar.ProgressBar
	writer  io.Writer
	enabled bool
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int, description string, enabled bool) *ProgressBar {
	if !enabled || total == 0 {
		return &ProgressBar{
			enabled: false,
			writer:  os.Stdout,
		}
	}

	bar := progressbar.NewOptions(total,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowIts(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
	)

	return &ProgressBar{
		bar:     bar,
		enabled: true,
		writer:  os.Stdout,
	}
}

// Increment advances the progress bar by one
func (p *ProgressBar) Increment() error {
	if !p.enabled || p.bar == nil {
		return nil
	}
	return p.bar.Add(1)
}

// Set sets the progress bar to a specific value
func (p *ProgressBar) Set(n int) error {
	if !p.enabled || p.bar == nil {
		return nil
	}
	return p.bar.Set(n)
}

// Finish completes the progress bar
func (p *ProgressBar) Finish() error {
	if !p.enabled || p.bar == nil {
		return nil
	}
	return p.bar.Finish()
}

// Clear clears the progress bar
func (p *ProgressBar) Clear() error {
	if !p.enabled || p.bar == nil {
		return nil
	}
	return p.bar.Clear()
}

// Printf prints a message (handles progress bar visibility)
func (p *ProgressBar) Printf(format string, args ...interface{}) {
	if p.enabled && p.bar != nil {
		p.bar.Clear()
		fmt.Fprintf(p.writer, format, args...)
		p.bar.RenderBlank()
	} else {
		fmt.Fprintf(p.writer, format, args...)
	}
}

// Println prints a message with newline
func (p *ProgressBar) Println(args ...interface{}) {
	if p.enabled && p.bar != nil {
		p.bar.Clear()
		fmt.Fprintln(p.writer, args...)
		p.bar.RenderBlank()
	} else {
		fmt.Fprintln(p.writer, args...)
	}
}

// Summary displays conversion statistics
func PrintSummary(total, success, failed, skipped int, duration time.Duration) {
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Conversion Summary")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Total:     %d files\n", total)
	fmt.Printf("Success:   %d files\n", success)
	if failed > 0 {
		fmt.Printf("Failed:    %d files\n", failed)
	}
	if skipped > 0 {
		fmt.Printf("Skipped:   %d files\n", skipped)
	}
	fmt.Printf("Duration:  %s\n", duration.Round(time.Millisecond))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// PrintError prints an error message
func PrintError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
}

// PrintWarning prints a warning message
func PrintWarning(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "WARNING: "+format+"\n", args...)
}

// PrintInfo prints an info message
func PrintInfo(format string, args ...interface{}) {
	fmt.Printf("INFO: "+format+"\n", args...)
}

// PrintVerbose prints a verbose message if verbose mode is enabled
func PrintVerbose(verbose bool, format string, args ...interface{}) {
	if verbose {
		fmt.Printf("VERBOSE: "+format+"\n", args...)
	}
}
