# ADR-002: Cobra & Viper Choice for CLI and Configuration

**Status**: Accepted
**Date**: 2025-10-17
**Deciders**: Development Team
**Technical Story**: CLI framework and configuration management

## Context and Problem Statement

SB requires a robust CLI framework and configuration management system. The system must support:
* Subcommands for different converters
* Global and command-specific flags
* Configuration files (YAML)
* Environment variables
* Flag precedence and overrides

Which frameworks should we use for CLI and configuration?

## Decision Drivers

* Maturity and community adoption
* Feature completeness
* Documentation quality
* Performance overhead
* Maintenance burden
* Integration between CLI and config

## Considered Options

### CLI Frameworks
1. **Cobra** - Full-featured CLI framework
2. **urfave/cli** - Simpler, less opinionated
3. **flag** (stdlib) - Minimal, standard library
4. **kingpin** - Alternative full-featured framework

### Configuration Libraries
1. **Viper** - Configuration with multiple sources
2. **envconfig** - Environment variables only
3. **toml/yaml** (stdlib-based) - Manual parsing
4. **koanf** - Alternative to Viper

## Decision Outcome

Chosen option: **Cobra + Viper**, because:

* Cobra and Viper integrate seamlessly (same author)
* Industry-standard (used by kubectl, Hugo, etc.)
* Excellent documentation and examples
* Auto-generated help and shell completions
* Viper's precedence system matches requirements
* Active maintenance and large community

### Positive Consequences

* Professional CLI UX out of the box
* Configuration hierarchy (flags > env > config > defaults)
* Auto-generated documentation
* Shell completion support (bash, zsh, fish)
* Familiar to Go developers
* Rich flag types and validation

### Negative Consequences

* Adds ~2MB to binary size
* Learning curve for contributors unfamiliar with Cobra
* Some "magic" behavior (automatic help, etc.)
* Potential over-engineering for simple CLIs

## Implementation

### Cobra Structure

```
cmd/
├── root.go           # Root command
├── formats/
│   └── mp4.go       # MP4 subcommand
├── ls.go            # List converters
├── info.go          # Media info
└── version.go       # Version info
```

### Viper Integration

```go
// Root command init
func init() {
    // Bind flags to Viper
    viper.BindPFlag("workers", rootCmd.PersistentFlags().Lookup("workers"))
    viper.BindPFlag("output_dir", rootCmd.PersistentFlags().Lookup("out"))

    // Environment variables
    viper.SetEnvPrefix("SB")
    viper.AutomaticEnv()

    // Config file
    viper.SetConfigName(".sb")
    viper.AddConfigPath("$HOME")
    viper.ReadInConfig()
}
```

### Configuration Precedence

1. CLI flags (highest priority)
2. Environment variables (SB_*)
3. Config file (~/.sb.yaml)
4. Defaults (lowest priority)

## Pros and Cons of Other Options

### urfave/cli
* **Good**: Simpler API, smaller footprint
* **Bad**: Less feature-rich, no integrated config solution
* **Outcome**: Rejected - insufficient for complex CLI

### flag (stdlib)
* **Good**: Zero dependencies, minimal
* **Bad**: Manual subcommand handling, no config integration
* **Outcome**: Rejected - too much custom code required

### kingpin
* **Good**: Feature-rich, good validation
* **Bad**: Less popular, less documentation
* **Outcome**: Rejected - Cobra more established

### envconfig
* **Good**: Simple, focused
* **Bad**: Env vars only, no YAML/flags
* **Outcome**: Rejected - need multi-source config

### koanf
* **Good**: More flexible than Viper
* **Bad**: No CLI integration, newer/less proven
* **Outcome**: Rejected - Viper + Cobra synergy too strong

## Examples

### Command Structure
```bash
sb mp4 --help                    # Auto-generated help
sb mp4 -w 8 -q 20 input.mov      # CLI flags
SB_WORKERS=8 sb mp4 input.mov    # Environment variables
sb mp4 input.mov                 # Config file settings
```

### Configuration File
```yaml
# ~/.sb.yaml
workers: 4
skip_existing: true
mp4:
  quality: 23
  preset: medium
```

## Performance Impact

* Binary size: +~2MB (Cobra + Viper)
* Startup time: +~10ms (negligible)
* Runtime overhead: Minimal (config loaded once)

## Links

* [Cobra Documentation](https://github.com/spf13/cobra)
* [Viper Documentation](https://github.com/spf13/viper)
* [Root Command Implementation](../../cmd/root.go)
* [Config Package](../../internal/config/config.go)
* Related ADRs: ADR-001 (Converter Interface)
