# ADR-001: Converter Interface Pattern

**Status**: Accepted
**Date**: 2025-10-17
**Deciders**: Development Team
**Technical Story**: Core architecture design for SB

## Context and Problem Statement

SB needs to support multiple media conversion formats (MP4, JPG, PNG, etc.) with a design that allows easy addition of new converters without modifying core code. How should we structure the converter system to maximize extensibility while maintaining simplicity?

## Decision Drivers

* Extensibility: Easy to add new converters
* Maintainability: Converters should be self-contained
* Discoverability: Automatic registration of converters
* Type Safety: Compile-time interface checking
* Performance: Minimal overhead from abstraction

## Considered Options

1. **Interface-based with Registry** - Define Converter interface, auto-register in init()
2. **Factory Pattern** - Manual registration with factory functions
3. **Plugin System** - External binaries loaded at runtime
4. **Hardcoded Switch** - Manual switch statement for each converter

## Decision Outcome

Chosen option: **Interface-based with Registry**, because:

* Provides compile-time type safety
* Enables automatic discovery via init() functions
* Self-contained converter implementations
* Zero external dependencies for core functionality
* Clean separation of concerns

### Positive Consequences

* New converters require minimal boilerplate
* Converters can be developed independently
* Easy to test converters in isolation
* Registry provides centralized access
* `sb ls` command can auto-discover converters

### Negative Consequences

* All converters compiled into binary (larger size)
* Cannot load converters dynamically at runtime
* init() ordering could cause subtle bugs (mitigated by design)

## Implementation

### Converter Interface

```go
type Converter interface {
    Name() string
    Description() string
    SupportedInputs() []string
    OutputExtension() string
    Validate(input string) error
    Convert(input string, opts Options) (*Result, error)
    ConvertBatch(inputs []string, opts Options) ([]*Result, error)
}
```

### Registry Pattern

```go
// Registry manages converters
var registry = &Registry{
    converters: make(map[string]Converter),
}

// Auto-register in converter package
func init() {
    converter.Register(NewMP4Converter())
}
```

### Usage

```go
// Get converter by name
conv, err := converter.Get("mp4")

// List all converters
converters := converter.ListConverters()
```

## Pros and Cons of Other Options

### Factory Pattern
* **Good**: Explicit registration, clear dependencies
* **Bad**: Manual registration required, more boilerplate
* **Outcome**: Rejected - too much manual work for contributors

### Plugin System
* **Good**: Smaller binary, runtime loading
* **Bad**: Complex versioning, security concerns, platform-specific
* **Outcome**: Rejected - unnecessary complexity for v1

### Hardcoded Switch
* **Good**: Simple, no abstraction
* **Bad**: Requires core changes for new converters, not extensible
* **Outcome**: Rejected - violates Open/Closed Principle

## Links

* [Internal Converter Package](../../internal/converter/)
* [MP4 Converter Implementation](../../internal/processors/mov_to_mp4/)
* [Registry Implementation](../../internal/converter/registry.go)
* Related ADRs: ADR-003 (Worker Pool Architecture)
