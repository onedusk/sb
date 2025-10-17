# Converter Development Guide

This guide explains how to add new converters to SB.

## Overview

A converter transforms media from one format to another. Converters implement a standard interface and self-register with the system, making them discoverable via `sb ls`.

## Quick Start

```bash
# 1. Create processor directory
mkdir -p internal/processors/heic_to_jpg

# 2. Create files
touch internal/processors/heic_to_jpg/processor.go
touch internal/processors/heic_to_jpg/options.go

# 3. Implement interface (see below)
# 4. Create CLI command
touch cmd/formats/jpg.go

# 5. Register command in main.go
# 6. Test
go test ./internal/processors/heic_to_jpg
make build
./dist/sb ls  # Verify converter appears
```

## Converter Interface

All converters must implement:

```go
type Converter interface {
    Name() string                    // Unique name (e.g., "mp4", "jpg")
    Description() string             // Human-readable description
    SupportedInputs() []string       // Input extensions (e.g., [".mov", ".avi"])
    OutputExtension() string         // Output extension (e.g., ".mp4")
    Validate(input string) error     // Validate input file
    Convert(input string, opts Options) (*Result, error)           // Single file
    ConvertBatch(inputs []string, opts Options) ([]*Result, error) // Batch processing
}
```

## Step-by-Step Implementation

### 1. Create Options Struct

**File**: `internal/processors/heic_to_jpg/options.go`

```go
package heic_to_jpg

type JPGOptions struct {
    Quality int    // JPEG quality (1-100)
    // Add converter-specific options
}

func DefaultJPGOptions() JPGOptions {
    return JPGOptions{
        Quality: 95,
    }
}

func (o *JPGOptions) Validate() error {
    if o.Quality < 1 || o.Quality > 100 {
        o.Quality = 95
    }
    return nil
}
```

### 2. Implement Converter

**File**: `internal/processors/heic_to_jpg/processor.go`

```go
package heic_to_jpg

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
    converter.Register(NewJPGConverter())
}

type JPGConverter struct {
    ffmpeg  *executor.FFmpeg
    options JPGOptions
}

func NewJPGConverter() *JPGConverter {
    return &JPGConverter{
        options: DefaultJPGOptions(),
    }
}

func (c *JPGConverter) Name() string {
    return "jpg"
}

func (c *JPGConverter) Description() string {
    return "Convert HEIC images to JPG format"
}

func (c *JPGConverter) SupportedInputs() []string {
    return []string{".heic", ".heif"}
}

func (c *JPGConverter) OutputExtension() string {
    return ".jpg"
}

func (c *JPGConverter) Validate(input string) error {
    // Check file exists
    info, err := os.Stat(input)
    if err != nil {
        return fmt.Errorf("cannot access file: %w", err)
    }

    if info.IsDir() {
        return fmt.Errorf("input is a directory, not a file")
    }

    // Check extension
    ext := strings.ToLower(filepath.Ext(input))
    for _, validExt := range c.SupportedInputs() {
        if ext == validExt {
            return nil
        }
    }

    return fmt.Errorf("unsupported file format: %s", ext)
}

func (c *JPGConverter) Convert(input string, opts converter.Options) (*converter.Result, error) {
    result := &converter.Result{Input: input}
    start := time.Now()

    // Validate
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

    // Skip if exists
    if opts.SkipExisting {
        if _, err := os.Stat(output); err == nil {
            result.Skipped = true
            result.SkipReason = "file already exists"
            result.Duration = time.Since(start)
            return result, nil
        }
    }

    // Dry run
    if opts.DryRun {
        fmt.Printf("[DRY-RUN] Would convert: %s -> %s\n", input, output)
        result.Success = true
        result.Duration = time.Since(start)
        return result, nil
    }

    // Create output directory
    if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
        result.Error = fmt.Errorf("failed to create output directory: %w", err)
        result.Duration = time.Since(start)
        return result, result.Error
    }

    // Build ffmpeg options
    ffmpegOpts := executor.FFmpegOptions{
        // Customize for your converter
        ExtraArgs: []string{"-q:v", fmt.Sprintf("%d", 100-c.options.Quality)},
        Verbose:   opts.Verbose,
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
        return result, result.Error
    }

    result.Success = true
    ui.PrintVerbose(opts.Verbose, "Successfully converted %s", input)

    return result, nil
}

func (c *JPGConverter) ConvertBatch(inputs []string, opts converter.Options) ([]*converter.Result, error) {
    if len(inputs) == 0 {
        return nil, fmt.Errorf("no input files provided")
    }

    results := make([]*converter.Result, 0, len(inputs))
    progress := ui.NewProgressBar(len(inputs), "Converting", opts.ShowProgress && !opts.DryRun)

    pool := executor.NewPool(opts.Workers)
    pool.Start()

    go func() {
        for _, input := range inputs {
            inputCopy := input
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

    errors := make([]error, 0)
    for err := range pool.Results() {
        if err != nil {
            errors = append(errors, err)
        }
    }

    progress.Finish()

    if len(errors) > 0 {
        return results, fmt.Errorf("%d conversion(s) failed", len(errors))
    }

    return results, nil
}

func (c *JPGConverter) SetOptions(opts JPGOptions) error {
    if err := opts.Validate(); err != nil {
        return err
    }
    c.options = opts
    return nil
}

func (c *JPGConverter) determineOutputPath(input string, opts converter.Options) string {
    baseName := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
    outputName := baseName + c.OutputExtension()

    if opts.OutputDir != "" {
        return filepath.Join(opts.OutputDir, outputName)
    }

    return filepath.Join(filepath.Dir(input), outputName)
}
```

### 3. Create CLI Command

**File**: `cmd/formats/jpg.go`

```go
package formats

import (
    "context"
    "fmt"

    "github.com/onedusk/sb/internal/config"
    "github.com/onedusk/sb/internal/converter"
    "github.com/onedusk/sb/internal/processors/heic_to_jpg"
    "github.com/spf13/cobra"
)

var (
    jpgQuality int
)

var JPGCmd = &cobra.Command{
    Use:   "jpg [files...]",
    Short: "Convert HEIC images to JPG format",
    Long: `Convert HEIC/HEIF images to JPG format.

Examples:
  sb jpg photo.heic
  sb jpg *.heic
  sb jpg -q 95 -o ./converted *.heic`,
    RunE: runJPGConvert,
}

func init() {
    JPGCmd.Flags().IntVarP(&jpgQuality, "quality", "q", 0, "JPEG quality (1-100, default: 95)")
}

func runJPGConvert(cmd *cobra.Command, args []string) error {
    cfg := config.Get()

    inputs, err := gatherInputs(args, "", false)
    if err != nil {
        return err
    }

    if len(inputs) == 0 {
        return fmt.Errorf("no input files found")
    }

    conv, err := converter.Get("jpg")
    if err != nil {
        return fmt.Errorf("jpg converter not available: %w", err)
    }

    jpgConv, ok := conv.(*heic_to_jpg.JPGConverter)
    if !ok {
        return fmt.Errorf("invalid converter type")
    }

    opts := heic_to_jpg.DefaultJPGOptions()
    if jpgQuality > 0 {
        opts.Quality = jpgQuality
    }

    if err := jpgConv.SetOptions(opts); err != nil {
        return err
    }

    convOpts := converter.Options{
        OutputDir:    cfg.OutputDir,
        Workers:      cfg.Workers,
        SkipExisting: cfg.SkipExisting,
        Verbose:      cfg.Verbose,
        ShowProgress: true,
        Context:      context.Background(),
    }

    if len(inputs) == 1 {
        _, err = jpgConv.Convert(inputs[0], convOpts)
        return err
    }

    _, err = jpgConv.ConvertBatch(inputs, convOpts)
    return err
}
```

### 4. Register Command

**File**: `main.go`

```go
import (
    "github.com/onedusk/sb/cmd/formats"
    _ "github.com/onedusk/sb/internal/processors/heic_to_jpg"  // Register
)

func main() {
    cmd.GetRootCmd().AddCommand(formats.JPGCmd)
    cmd.Execute()
}
```

### 5. Add Tests

**File**: `internal/processors/heic_to_jpg/processor_test.go`

```go
package heic_to_jpg

import (
    "testing"
)

func TestJPGConverter_Name(t *testing.T) {
    conv := NewJPGConverter()
    if conv.Name() != "jpg" {
        t.Errorf("Name() = %v, want %v", conv.Name(), "jpg")
    }
}

func TestJPGConverter_Validate(t *testing.T) {
    conv := NewJPGConverter()

    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid HEIC", "test.heic", false},
        {"invalid ext", "test.txt", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := conv.Validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Best Practices

1. **Reuse Worker Pool**: Use provided executor.Pool for batch processing
2. **Support Dry Run**: Always check opts.DryRun
3. **Respect SkipExisting**: Check before overwriting
4. **Verbose Logging**: Use ui.PrintVerbose
5. **Error Wrapping**: Use fmt.Errorf with %w
6. **Context Cancellation**: Check context in long operations
7. **Progress Tracking**: Use ui.NewProgressBar
8. **Input Validation**: Validate early in Validate()

## Testing Checklist

- [ ] Unit tests for Validate()
- [ ] Test with actual conversion
- [ ] Test batch processing
- [ ] Test dry-run mode
- [ ] Test skip existing
- [ ] Test with various file types
- [ ] Test error cases

## Documentation Checklist

- [ ] Update README.md with new converter
- [ ] Add command examples
- [ ] Update CHANGELOG.md
- [ ] Add converter to docs/architecture.md
- [ ] Document converter-specific options

## See Also

- [MP4 Converter](../../internal/processors/mov_to_mp4/) - Reference implementation
- [Converter Interface](../../internal/converter/converter.go)
- [Executor Package](../../internal/executor/)
- [Contributing Guide](../CONTRIBUTING.md)
