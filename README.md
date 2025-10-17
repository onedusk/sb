# SB - Scalable Media Processor

A powerful, extensible CLI tool for media processing with support for batch conversion, hardware acceleration, and quality controls.

## Features

- **Multiple Format Support**: Convert between various video and image formats
- **Batch Processing**: Process multiple files in parallel with configurable worker pools
- **Hardware Acceleration**: Support for VideoToolbox (macOS), NVENC (NVIDIA), and QSV (Intel)
- **Quality Controls**: Fine-tune output with CRF, presets, bitrate, and codec options
- **Flexible Configuration**: CLI flags, config files, and environment variables
- **Progress Tracking**: Real-time progress bars and conversion statistics
- **Extensible Architecture**: Plugin-like converter system for easy expansion
- **Single Binary**: Distributed as a standalone binary with no runtime dependencies (except ffmpeg)

## Installation

### Prerequisites

- **ffmpeg** must be installed and available in your PATH
  - macOS: `brew install ffmpeg`
  - Linux: `apt install ffmpeg` or `yum install ffmpeg`
  - Windows: Download from [ffmpeg.org](https://ffmpeg.org)

### Building from Source

```bash
# Clone and build
cd scripts/media/sb
make build

# Install to /usr/local/bin
make install

# Or build for all platforms
make release
```

### Pre-built Binaries

Download pre-built binaries from the releases page (darwin-amd64, darwin-arm64, linux-amd64, linux-arm64, windows-amd64).

## Quick Start

```bash
# Convert a single MOV file to MP4
sb mp4 video.mov

# Convert all MOV files in current directory
sb mp4 *.mov

# Convert directory recursively with 8 workers
sb mp4 -d ./videos -r -w 8

# High quality conversion with hardware acceleration
sb mp4 -q 18 -p slow --hw videotoolbox input.mov

# Batch convert with custom output directory
sb mp4 -o ./converted -s *.mov
```

## Usage

### Global Flags

```
-w, --workers N      Number of parallel workers (default: CPU count)
-o, --out DIR        Output directory
-r, --recursive      Process directories recursively
-s, --skip           Skip existing files
-n, --dry-run        Preview without converting
-v, --verbose        Verbose output
-f, --flat           Flatten output directory structure
    --config FILE    Config file (default: $HOME/.sb.yaml)
```

### MP4 Conversion

Convert video files to MP4 format using H.264/H.265 encoding.

```bash
sb mp4 [files...] [flags]
```

**Supported Input Formats**: .mov, .avi, .mkv, .flv, .wmv, .m4v, .mpeg, .mpg, .webm

**MP4-Specific Flags:**

```
-q, --quality N           CRF quality (0-51, lower = better, default: 23)
-p, --preset PRESET       Encoding preset (ultrafast|fast|medium|slow|veryslow)
-c, --codec CODEC         Video codec (h264|h265|vp9)
    --audio CODEC         Audio codec (aac|mp3|copy)
    --audio-bitrate RATE  Audio bitrate (e.g., 128k, 192k)
-b, --bitrate RATE        Video bitrate (e.g., 2M, 5M)
    --hw TYPE             Hardware acceleration (videotoolbox|nvenc|qsv)
-d, --dir DIR             Input directory
    --recursive           Process directory recursively
```

**Examples:**

```bash
# Basic conversion
sb mp4 video.mov

# High quality, slow encoding
sb mp4 -q 18 -p slow video.mov

# Hardware accelerated (macOS)
sb mp4 --hw videotoolbox --codec h265 *.mov

# Batch convert with custom settings
sb mp4 -q 20 -w 4 -o ./converted *.mov

# Directory processing
sb mp4 -d ~/Videos -r -s -w 8

# Bitrate control instead of CRF
sb mp4 -b 5M --audio-bitrate 192k video.mov
```

### Utility Commands

```bash
# List available converters
sb ls

# Show media file info
sb info video.mp4

# Show version
sb version
```

## Configuration

SB supports configuration files for setting defaults.

### Configuration File

Create `~/.sb.yaml` or `.sb.yaml` in your project:

```yaml
# Global settings
workers: 4
skip_existing: true
output_dir: "./converted"
flat_structure: false
verbose: false

# MP4 conversion settings
mp4:
  quality: 23
  preset: medium
  codec: h264
  audio: aac
  bitrate: ""
  hardware:
    enabled: false
    type: videotoolbox
```

### Environment Variables

All settings can be overridden with environment variables prefixed with `SB_`:

```bash
export SB_WORKERS=8
export SB_MP4_QUALITY=20
export SB_MP4_PRESET=slow
```

### Priority Order

1. CLI flags (highest)
2. Environment variables
3. Config file
4. Defaults (lowest)

## Architecture

SB uses a plugin-like converter architecture for extensibility:

```
sb/
├── cmd/                    # Cobra CLI commands
│   ├── root.go            # Root command
│   ├── formats/           # Format commands (mp4, jpg, etc.)
│   └── version.go         # Utility commands
├── internal/
│   ├── converter/         # Converter interface & registry
│   ├── processors/        # Converter implementations
│   │   └── mov_to_mp4/   # MP4 converter
│   ├── executor/          # FFmpeg wrapper & worker pool
│   ├── config/            # Viper configuration
│   └── ui/                # Progress bars & output
└── main.go
```

### Adding New Converters

1. Create processor in `internal/processors/`
2. Implement `converter.Converter` interface
3. Register in `init()` function
4. Add command in `cmd/formats/`
5. Register command in `main.go`

## Examples

### Basic Workflow

```bash
# Convert single file with defaults
sb mp4 input.mov
# Output: input.mp4 (in same directory)

# Convert with custom quality
sb mp4 -q 18 input.mov
# Higher quality (lower CRF = better quality)

# Convert multiple files
sb mp4 video1.mov video2.mov video3.mov
```

### Batch Processing

```bash
# Convert all MOV files in directory
sb mp4 *.mov

# Recursive directory conversion
sb mp4 -d ~/Videos -r -w 8 -o ~/Converted

# Skip existing files (resume conversion)
sb mp4 -d ~/Videos -r -s
```

### Quality Optimization

```bash
# Maximum quality (slow)
sb mp4 -q 18 -p veryslow input.mov

# Balanced quality/speed
sb mp4 -q 23 -p medium input.mov

# Fast conversion (lower quality)
sb mp4 -q 28 -p fast input.mov

# Bitrate control
sb mp4 -b 2M input.mov
```

### Hardware Acceleration

```bash
# macOS VideoToolbox
sb mp4 --hw videotoolbox -c h264 *.mov

# NVIDIA NVENC (Linux/Windows)
sb mp4 --hw nvenc -c h264 *.mov

# Intel QSV
sb mp4 --hw qsv -c h264 *.mov
```

### Advanced Usage

```bash
# Dry run to preview
sb mp4 -n -v *.mov

# Verbose output for debugging
sb mp4 -v input.mov

# Custom output structure
sb mp4 -o ./converted -f *.mov
# (flat structure: all files in output dir)

# Using config file
sb mp4 --config ./custom.yaml *.mov
```

## Troubleshooting

### ffmpeg not found

```bash
# Install ffmpeg
brew install ffmpeg  # macOS
apt install ffmpeg   # Linux
```

### Hardware acceleration not working

```bash
# Check available encoders
ffmpeg -encoders | grep videotoolbox  # macOS
ffmpeg -encoders | grep nvenc          # NVIDIA
ffmpeg -encoders | grep qsv            # Intel

# Test hardware encoder
sb mp4 --hw videotoolbox -v test.mov
```

### Conversion fails

```bash
# Enable verbose output
sb mp4 -v input.mov

# Check input file
sb info input.mov

# Try without hardware acceleration
sb mp4 input.mov
```

## Performance Tips

1. **Use hardware acceleration** for faster encoding (2-5x speedup)
2. **Adjust worker count** based on CPU cores and workload
3. **Balance quality and speed** with preset selection
4. **Use CRF** for consistent quality, bitrate for size control
5. **Skip existing files** when resuming large batch jobs

## Contributing

This tool is part of the Dusk Labs monorepo. Contributions welcome!

### Development

```bash
# Run tests
make test

# Build for development
make build

# Run with changes
make run

# Build for all platforms
make release
```

## License

See repository root for license information.

## Version

Check version with:

```bash
sb version
```
