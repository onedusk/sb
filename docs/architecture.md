# SB Architecture

This document describes the architecture and design of the SB media processor.

## Overview

SB is a command-line media processing toolkit built in Go with a plugin-like architecture for converters. It wraps ffmpeg for media conversion while providing parallelization, configuration management, and a user-friendly CLI.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI Layer                            │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │ sb mp4   │  │  sb ls   │  │ sb info  │  │sb version│   │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘   │
│       │             │               │             │          │
│       └─────────────┴───────────────┴─────────────┘          │
│                         │                                     │
│                    Cobra Framework                            │
└──────────────────────────┬────────────────────────────────────┘
                           │
┌──────────────────────────┴────────────────────────────────────┐
│                   Core Components                              │
│  ┌────────────────┐    ┌────────────────┐                    │
│  │  Converter     │    │  Config        │                    │
│  │  Registry      │◄───┤  Manager       │                    │
│  │                │    │  (Viper)       │                    │
│  └────────┬───────┘    └────────────────┘                    │
│           │                                                    │
│  ┌────────▼──────────────────────────────────────┐           │
│  │  Converter Interface                           │           │
│  │  - Name(), Description()                       │           │
│  │  - SupportedInputs(), OutputExtension()        │           │
│  │  - Validate(input)                             │           │
│  │  - Convert(input, opts)                        │           │
│  │  - ConvertBatch(inputs, opts)                  │           │
│  └────────┬──────────────────────────────────────┘           │
└───────────┼──────────────────────────────────────────────────┘
            │
┌───────────┴──────────────────────────────────────────────────┐
│                 Processor Layer                               │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐            │
│  │  MP4       │  │  HEIC      │  │   Future   │            │
│  │ Processor  │  │ Processor  │  │ Processors │            │
│  └─────┬──────┘  └────────────┘  └────────────┘            │
│        │                                                      │
│        ▼                                                      │
│  ┌─────────────────────┐                                     │
│  │  Executor Layer     │                                     │
│  │  ┌────────────┐     │                                     │
│  │  │ FFmpeg     │     │                                     │
│  │  │ Wrapper    │     │                                     │
│  │  └─────┬──────┘     │                                     │
│  │  ┌─────▼──────┐     │                                     │
│  │  │ Worker     │     │                                     │
│  │  │ Pool       │     │                                     │
│  │  └────────────┘     │                                     │
│  └─────────────────────┘                                     │
└──────────────┬───────────────────────────────────────────────┘
               │
┌──────────────▼───────────────────────────────────────────────┐
│                 External Dependencies                         │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐            │
│  │  ffmpeg    │  │  ffprobe   │  │    OS      │            │
│  │            │  │            │  │  (exec)    │            │
│  └────────────┘  └────────────┘  └────────────┘            │
└───────────────────────────────────────────────────────────────┘
```

## Component Descriptions

### CLI Layer (`cmd/`)

Built with Cobra, provides user interface:
- **Root Command**: Global flags, initialization
- **Format Commands**: Converter-specific subcommands (mp4, jpg, etc.)
- **Utility Commands**: ls, info, version

### Converter System (`internal/converter/`)

Plugin-like architecture:
- **Interface**: Defines converter contract
- **Registry**: Auto-discovery and management
- **Types**: Shared types (Options, Result, Stats)

### Processors (`internal/processors/`)

Converter implementations:
- Each processor in own package
- Auto-registers in init()
- Self-contained with options validation

### Executor (`internal/executor/`)

Execution infrastructure:
- **FFmpeg Wrapper**: Command building and execution
- **Worker Pool**: Parallel job processing
- Context-based cancellation

### Configuration (`internal/config/`)

Multi-source configuration:
- Viper-based management
- Precedence: CLI > Env > Config > Defaults
- Type-safe unmarshaling

### UI (`internal/ui/`)

User interface components:
- Progress bars
- Statistics summaries
- Formatted output

## Data Flow

### Single File Conversion

```
User Command
    ↓
CLI Parser (Cobra)
    ↓
Get Converter (Registry)
    ↓
Load Config (Viper)
    ↓
Validate Input
    ↓
Build FFmpeg Options
    ↓
Execute FFmpeg
    ↓
Collect Result
    ↓
Display Output
```

### Batch Conversion

```
User Command (multiple files)
    ↓
CLI Parser
    ↓
Get Converter
    ↓
Create Worker Pool (N workers)
    ↓
Submit Jobs (one per file)
    ↓
Workers Process Jobs in Parallel
    ├─ Worker 1 → FFmpeg → Result
    ├─ Worker 2 → FFmpeg → Result
    ├─ Worker 3 → FFmpeg → Result
    └─ Worker N → FFmpeg → Result
    ↓
Aggregate Results
    ↓
Display Summary
```

## Extension Points

### Adding a New Converter

1. **Create Processor Package**
   ```
   internal/processors/new_format/
   ├── processor.go
   └── options.go
   ```

2. **Implement Interface**
   ```go
   type NewConverter struct {}
   func (c *NewConverter) Name() string { ... }
   // ... implement all interface methods
   ```

3. **Register Converter**
   ```go
   func init() {
       converter.Register(NewNewConverter())
   }
   ```

4. **Create CLI Command**
   ```
   cmd/formats/newformat.go
   ```

5. **Register Command**
   ```go
   // main.go
   cmd.GetRootCmd().AddCommand(formats.NewFormatCmd)
   ```

## Design Patterns

### Registry Pattern
Converters self-register using init() functions, enabling auto-discovery.

### Strategy Pattern
Converter interface allows swapping conversion strategies.

### Factory Pattern
Registry acts as factory for converter instantiation.

### Worker Pool Pattern
Fixed pool of workers process jobs from queue.

### Command Pattern
Jobs encapsulate conversion operations.

## Concurrency Model

- **Worker Pool**: Fixed number of goroutines
- **Job Queue**: Channel-based work distribution
- **Context Propagation**: Cancellation signal
- **Result Collection**: Aggregation channel
- **Synchronization**: WaitGroup for completion

## Error Handling

- Converters return errors, don't panic
- Errors wrapped with context
- Batch operations collect all errors
- Non-fatal errors don't stop other jobs
- Clear error messages for users

## Security Considerations

See [SECURITY.md](../SECURITY.md) for full details:
- Command injection prevention (args as array)
- Input validation (file paths, extensions)
- No elevated privileges required
- FFmpeg security inherited

## Performance Characteristics

- **Startup**: < 100ms
- **Memory**: Base ~10MB + workers
- **Parallelization**: 4-8x speedup typical
- **Overhead**: < 1% for batch operations

## Testing Strategy

- Unit tests for core logic
- Integration tests with real ffmpeg
- Converter interface compliance tests
- Mock ffmpeg for unit tests
- CI runs all tests on PR

## Configuration Architecture

```
Priority (High to Low):
1. CLI Flags:        --workers=8
2. Env Variables:    SB_WORKERS=8
3. Config File:      ~/.sb.yaml
4. Defaults:         runtime.NumCPU()

Viper merges all sources automatically.
```

## Deployment

- Single static binary
- No runtime dependencies (except ffmpeg)
- Cross-platform builds via Makefile
- GitHub Actions for CI/CD

## Future Architecture

- Distributed processing coordinator
- Plugin system for external converters
- RESTful API mode
- Streaming/pipeline mode

## Related Documents

- [ADR Directory](adr/) - Architecture decisions
- [PRD](PRD.md) - Product requirements
- [Development Guide](development.md) - Developer workflow
- [Converter Guide](converters.md) - Adding converters
