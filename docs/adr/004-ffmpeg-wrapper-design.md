# ADR-004: FFmpeg Wrapper Design

**Status**: Accepted
**Date**: 2025-10-17
**Deciders**: Development Team
**Technical Story**: FFmpeg integration strategy

## Context and Problem Statement

SB needs to leverage ffmpeg for media conversion. FFmpeg is a powerful, mature tool with extensive format and codec support. How should we integrate ffmpeg into SB? Should we use command-line execution or native library bindings?

## Decision Drivers

* Simplicity: Easy to implement and maintain
* Compatibility: Work with any ffmpeg version
* Performance: Minimize overhead
* Features: Access to full ffmpeg capabilities
* Error Handling: Clear error messages
* Security: Prevent command injection
* Portability: Cross-platform support

## Considered Options

1. **Command-Line Wrapper** - Execute ffmpeg via os/exec
2. **CGo Bindings** - Link against libav* libraries
3. **Pure Go Implementation** - Rewrite conversion logic in Go
4. **REST API** - Call external ffmpeg service

## Decision Outcome

Chosen option: **Command-Line Wrapper**, because:

* Simple implementation with stdlib (os/exec)
* No build dependencies (users install ffmpeg separately)
* Works with any ffmpeg version/build
* Full access to ffmpeg features via flags
* Clear separation of concerns
* Cross-platform compatible
* No C build complexity

### Positive Consequences

* Simple Go code, easy to maintain
* Users can use their preferred ffmpeg build
* No CGo complexity or build issues
* Easy to add new ffmpeg features (just flags)
* Clear process boundaries for debugging
* ffmpeg updates don't require SB recompile

### Negative Consequences

* Requires ffmpeg to be installed separately
* Process spawn overhead (~10ms per file)
* No direct access to ffmpeg internals
* Error messages are ffmpeg stderr parsing
* Dependency on ffmpeg CLI interface stability

## Implementation

### FFmpeg Wrapper Structure

```go
type FFmpeg struct {
    binaryPath string
}

type FFmpegOptions struct {
    VideoCodec   string
    CRF          int
    Preset       string
    AudioCodec   string
    HWAccel      string
    // ... additional options
}

func (f *FFmpeg) Convert(ctx context.Context, input, output string, opts FFmpegOptions) (*FFmpegResult, error) {
    args := f.buildArgs(input, output, opts)
    cmd := exec.CommandContext(ctx, f.binaryPath, args...)
    // Execute and return result
}
```

### Argument Building

```go
func (f *FFmpeg) buildArgs(input, output string, opts FFmpegOptions) []string {
    args := []string{"-y"} // Overwrite output

    // Hardware acceleration (before input)
    if opts.HWAccel != "" {
        args = append(args, "-hwaccel", opts.HWAccel)
    }

    // Input
    args = append(args, "-i", input)

    // Video codec
    args = append(args, "-c:v", opts.VideoCodec)

    // Quality
    if opts.CRF > 0 {
        args = append(args, "-crf", fmt.Sprintf("%d", opts.CRF))
    }

    // ... build complete arg list

    args = append(args, output)
    return args
}
```

### Security: Preventing Command Injection

* Arguments passed as array, not string
* No shell interpolation
* Input paths validated before use
* User options validated against whitelist

```go
// Safe - args are array elements
cmd := exec.Command("ffmpeg", "-i", userInput, output)

// Unsafe - don't do this
cmd := exec.Command("sh", "-c", "ffmpeg -i " + userInput)
```

### Error Handling

```go
type FFmpegResult struct {
    Success  bool
    Duration time.Duration
    Stdout   string
    Stderr   string  // Parse for error details
    Error    error
}
```

## Pros and Cons of Other Options

### CGo Bindings (libav*)
* **Good**: No process spawn, direct API access, better performance
* **Bad**: Complex build, platform-specific, CGo overhead, version coupling
* **Outcome**: Rejected - build complexity outweighs benefits

### Pure Go Implementation
* **Good**: No dependencies, full control
* **Bad**: Massive effort, incomplete codec support, ongoing maintenance
* **Outcome**: Rejected - unrealistic scope

### REST API
* **Good**: Language-agnostic, scalable
* **Bad**: Network overhead, deployment complexity, not suited for CLI tool
* **Outcome**: Rejected - wrong architecture for use case

## Performance Analysis

### Process Spawn Overhead
* Spawn time: ~10ms per invocation
* For batch of 100 files: ~1 second overhead
* Conversion time: Seconds to minutes per file
* Overhead: < 1% of total time

### Memory Usage
* Each ffmpeg process: 50-200MB
* Multiple workers: Memory scales linearly
* SB wrapper: < 10MB

### Optimization Strategies
* Reuse FFmpeg struct (no benefit from pooling processes)
* Let ffmpeg handle parallelism internally (not implemented)
* Use hardware acceleration when available

## ffmpeg Version Compatibility

### Minimum Version
* Recommended: ffmpeg >= 4.0
* Fallback: Most features work on 3.x

### Version Detection
```go
func (f *FFmpeg) CheckVersion() (string, error) {
    cmd := exec.Command(f.binaryPath, "-version")
    output, err := cmd.Output()
    // Parse version from output
}
```

### Feature Detection
* Probe for hardware encoders: `ffmpeg -encoders | grep videotoolbox`
* Check codec support: `ffmpeg -codecs | grep h265`

## Alternative Executables

### ffprobe Integration
```go
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
    output, err := cmd.Output()
    return string(output), err
}
```

## Hardware Acceleration Support

### VideoToolbox (macOS)
```bash
sb mp4 --hw videotoolbox --codec h264 input.mov
# Translates to: ffmpeg -hwaccel videotoolbox -c:v h264_videotoolbox
```

### NVENC (NVIDIA)
```bash
sb mp4 --hw nvenc --codec h264 input.mov
# Translates to: ffmpeg -hwaccel nvenc -c:v h264_nvenc
```

## Future Enhancements

* **Progress Parsing**: Parse ffmpeg stderr for progress percentage
* **Quality Presets**: Map user-friendly names to ffmpeg options
* **Codec Auto-detection**: Choose best codec for hardware
* **Filter Chains**: Support complex ffmpeg filters
* **Multiple Passes**: Two-pass encoding for better quality

## Testing Strategy

### Unit Tests
* Mock ffmpeg with test binary
* Verify argument building
* Test error handling

### Integration Tests
* Requires ffmpeg installed
* Test actual conversions
* Verify output format

## Links

* [FFmpeg Wrapper Implementation](../../internal/executor/ffmpeg.go)
* [FFmpeg Documentation](https://ffmpeg.org/documentation.html)
* [Go os/exec Package](https://pkg.go.dev/os/exec)
* Related ADRs: ADR-003 (Worker Pool Architecture)

## References

* [FFmpeg CLI Documentation](https://ffmpeg.org/ffmpeg.html)
* [Hardware Acceleration Guide](https://trac.ffmpeg.org/wiki/HWAccelIntro)
* [os/exec Best Practices](https://go.dev/blog/context)
