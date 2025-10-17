# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Additional converter formats planned (HEIC→JPG, image resizing, audio extraction)
- Path preservation option for nested directory structures
- Configuration validation command
- Batch job templates

## [0.1.0] - 2025-10-17

### Added
- Initial release of SB media processor
- MP4 converter with MOV→MP4 support
- Multiple input format support (.mov, .avi, .mkv, .flv, .wmv, .m4v, .mpeg, .mpg, .webm)
- Parallel batch processing with configurable worker pools
- Hardware acceleration support (VideoToolbox, NVENC, QSV)
- Quality controls (CRF, preset, bitrate, codec selection)
- Viper-based configuration system with YAML support
- Environment variable configuration (SB_* prefix)
- CLI built with Cobra framework
- Progress tracking with real-time statistics
- Dry-run mode for previewing conversions
- Skip existing files option for resume capability
- Verbose logging mode
- Utility commands (ls, info, version)
- Cross-platform build support (darwin, linux, windows)
- Makefile for build automation
- Comprehensive README documentation
- Example configuration file (.sb.yaml.example)

### Architecture
- Converter interface and registry system for extensibility
- FFmpeg wrapper with option building
- Worker pool pattern for concurrent processing
- Modular internal package structure
- Plugin-like converter registration

[Unreleased]: https://github.com/dusk-labs/sb/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/dusk-labs/sb/releases/tag/v0.1.0
