# Product Requirements Document: SB Media Processor

**Version**: 1.0
**Last Updated**: 2025-10-17
**Status**: Active

## Executive Summary

SB (Scalable Builder/Batch) is a command-line media processing toolkit designed for developers, content creators, and automation workflows. It provides a unified interface for converting between media formats with emphasis on performance, extensibility, and ease of use.

## Vision & Goals

### Vision
Create a Go-based, extensible media processing toolkit that becomes the go-to CLI tool for batch media conversions, offering superior performance through parallelization and hardware acceleration.

### Primary Goals
1. **Simplicity**: Easy-to-use CLI with sensible defaults
2. **Performance**: Parallel processing with hardware acceleration support
3. **Extensibility**: Plugin-like architecture for adding new converters
4. **Reliability**: Robust error handling and resume capability
5. **Distribution**: Single binary with minimal dependencies

### Success Metrics
- Processing speed: 4-8x faster than sequential processing
- User adoption: GitHub stars, downloads, community contributions
- Reliability: < 1% failure rate on valid media files
- Extension: Community-contributed converters

## Target Users

### Primary Users
1. **Content Creators**
   - Photographers converting RAW/HEIC to standard formats
   - Video editors batch processing footage
   - Podcasters converting audio formats

2. **Developers**
   - Automation scripts for media pipelines
   - Build systems requiring media conversion
   - Server-side media processing

3. **System Administrators**
   - Batch migration between formats
   - Storage optimization workflows
   - Media library maintenance

### Use Cases

#### UC-001: Batch Video Conversion
**User**: Video editor
**Goal**: Convert 100+ MOV files from iPhone to MP4 for editing
**Flow**:
1. User runs: `sb mp4 -d ~/iPhone-Videos -r -w 8`
2. SB processes files in parallel using 8 workers
3. Progress bar shows real-time status
4. Completed files appear in same directory structure

**Success Criteria**: All files converted with < 2% failure rate, 5x faster than sequential

#### UC-002: Hardware-Accelerated Encoding
**User**: Developer on Mac
**Goal**: Use VideoToolbox for fast encoding
**Flow**:
1. User runs: `sb mp4 --hw videotoolbox -c h265 *.mov`
2. SB utilizes hardware encoder
3. Conversion completes 3-5x faster than software encoding

**Success Criteria**: Hardware acceleration utilized, significant speedup observed

#### UC-003: Automated CI/CD Pipeline
**User**: DevOps engineer
**Goal**: Convert media assets during build
**Flow**:
1. CI script runs: `sb mp4 -o ./build/assets -s -v input/*.mov`
2. SB skips already converted files
3. Logs output for build system
4. Exit code indicates success/failure

**Success Criteria**: Reliable automation, proper exit codes, resume capability

#### UC-004: Custom Quality Settings
**User**: Content creator
**Goal**: High-quality archival conversion
**Flow**:
1. User creates `~/.sb.yaml` with quality settings
2. Runs: `sb mp4 -q 18 -p veryslow important-footage.mov`
3. SB applies quality settings

**Success Criteria**: Fine-grained quality control, configuration persistence

## Features

### MVP (v0.1 - Current)

#### Core Features
- ✓ MP4 converter (MOV→MP4 primary, multi-format support)
- ✓ Parallel batch processing with worker pools
- ✓ Hardware acceleration (VideoToolbox, NVENC, QSV)
- ✓ Quality controls (CRF, preset, bitrate, codec)
- ✓ Configuration system (YAML, env vars, CLI flags)
- ✓ Progress tracking and statistics
- ✓ Dry-run mode
- ✓ Skip existing files (resume capability)
- ✓ Verbose logging

#### CLI Commands
- ✓ `sb mp4` - Convert to MP4
- ✓ `sb ls` - List converters
- ✓ `sb info` - Show media info
- ✓ `sb version` - Show version

#### Technical Features
- ✓ Single binary distribution
- ✓ Cross-platform (darwin, linux, windows)
- ✓ Extensible converter architecture
- ✓ FFmpeg wrapper with option building
- ✓ Context-based cancellation

### v0.2 (Planned)

#### New Converters
- [ ] HEIC → JPG converter
- [ ] PNG optimization
- [ ] Audio extraction (video → audio)
- [ ] GIF creation from video

#### Enhancements
- [ ] Path preservation for nested directories
- [ ] Configuration validation command: `sb config validate`
- [ ] Batch job templates
- [ ] Custom output filename patterns
- [ ] Media file filtering (size, duration, resolution)
- [ ] Progress resumption with job tracking

#### Quality of Life
- [ ] Interactive mode for settings
- [ ] Preset management (save/load quality presets)
- [ ] Conversion profiles (web, mobile, archival)
- [ ] Shell completions (bash, zsh, fish)

### v0.3 (Future)

#### Advanced Features
- [ ] Image resizing/scaling converter
- [ ] Watermarking support
- [ ] Subtitle extraction/embedding
- [ ] Batch audio normalization
- [ ] Multi-step pipelines (resize + convert + optimize)

#### Performance
- [ ] Distributed processing (multiple machines)
- [ ] GPU acceleration for applicable converters
- [ ] Caching for repeated operations
- [ ] Incremental conversions

#### Enterprise Features
- [ ] S3/cloud storage integration
- [ ] Webhook notifications
- [ ] Metrics export (Prometheus)
- [ ] RESTful API mode

## Technical Requirements

### Functional Requirements

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| FR-001 | Convert video files to MP4 format | P0 | ✓ Complete |
| FR-002 | Parallel processing with configurable workers | P0 | ✓ Complete |
| FR-003 | Hardware acceleration support | P0 | ✓ Complete |
| FR-004 | Quality control (CRF, preset, bitrate) | P0 | ✓ Complete |
| FR-005 | Configuration file support | P0 | ✓ Complete |
| FR-006 | Progress tracking | P0 | ✓ Complete |
| FR-007 | Skip existing files | P0 | ✓ Complete |
| FR-008 | Multiple input formats | P0 | ✓ Complete |
| FR-009 | Dry-run mode | P1 | ✓ Complete |
| FR-010 | Verbose logging | P1 | ✓ Complete |
| FR-011 | Additional converters (image, audio) | P1 | Planned |
| FR-012 | Path preservation | P2 | Planned |
| FR-013 | Job templates | P2 | Planned |

### Non-Functional Requirements

| ID | Requirement | Target | Status |
|----|-------------|--------|--------|
| NFR-001 | Performance: 4-8x speedup with parallelization | 4-8x | ✓ Met |
| NFR-002 | Binary size: < 20MB | < 20MB | ✓ Met (~10MB) |
| NFR-003 | Memory usage: < 100MB base + workers | Measured | ✓ Met |
| NFR-004 | Startup time: < 100ms | < 100ms | ✓ Met |
| NFR-005 | Error rate: < 1% on valid files | < 1% | Testing |
| NFR-006 | Code coverage: > 70% | > 70% | In progress |
| NFR-007 | Build time: < 30s | < 30s | ✓ Met |
| NFR-008 | Platform support: macOS, Linux, Windows | 3 platforms | ✓ Met |

### Dependencies

**External**:
- FFmpeg (user-installed, >= 4.0 recommended)

**Go Libraries**:
- github.com/spf13/cobra - CLI framework
- github.com/spf13/viper - Configuration
- github.com/schollz/progressbar/v3 - Progress UI
- Standard library (os/exec, sync, filepath, etc.)

**Development**:
- Go 1.21+
- golangci-lint
- make

## Architecture

See [architecture.md](architecture.md) for detailed architecture documentation.

### Key Design Decisions

See [adr/](adr/) directory for Architecture Decision Records:
- ADR-001: Converter Interface Pattern
- ADR-002: Cobra & Viper Choice
- ADR-003: Worker Pool Architecture
- ADR-004: FFmpeg Wrapper Design

## Roadmap

### Q4 2025
- ✓ v0.1.0: Initial release with MP4 converter
- [ ] v0.2.0: Additional converters (HEIC, PNG, audio)
- [ ] Community feedback and bug fixes
- [ ] Documentation improvements

### Q1 2026
- [ ] v0.3.0: Advanced features (pipelines, filtering)
- [ ] Performance optimizations
- [ ] Enterprise features exploration

### Q2 2026
- [ ] v1.0.0: Stable release
- [ ] API stabilization
- [ ] Long-term support commitment

## Out of Scope

The following are explicitly out of scope for SB:

1. **GUI Application**: SB is CLI-only
2. **Media Editing**: No editing capabilities (trimming, effects, etc.)
3. **Media Player**: No playback functionality
4. **FFmpeg Replacement**: SB wraps ffmpeg, doesn't replace it
5. **Cloud Service**: No hosted/SaaS offering
6. **Real-time Processing**: Focus is on batch operations
7. **DRM Handling**: No DRM circumvention or protected content

## Open Questions

1. **Distributed Processing**: How to coordinate multiple machines?
2. **Cloud Integration**: Which cloud providers to prioritize?
3. **Plugin System**: Should converters be pluggable via external binaries?
4. **Web UI**: Is there demand for a companion web interface?
5. **Licensing**: What license for the project?

## Appendix

### Glossary

- **Converter**: Module that transforms media from one format to another
- **Worker Pool**: Concurrent processing pattern with fixed worker count
- **CRF**: Constant Rate Factor, quality metric for video encoding
- **Hardware Acceleration**: Using GPU/dedicated hardware for encoding

### References

- FFmpeg documentation: https://ffmpeg.org/documentation.html
- Go CLI best practices: https://go.dev/doc/effective_go
- Cobra framework: https://github.com/spf13/cobra
- Viper config: https://github.com/spf13/viper

---

**Document History**:
- 2025-10-17: Initial PRD (v1.0)
