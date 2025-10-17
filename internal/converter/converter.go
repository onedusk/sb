package converter

// Converter defines the interface that all media converters must implement
type Converter interface {
	// Name returns the unique name of this converter
	Name() string

	// Description returns a human-readable description
	Description() string

	// SupportedInputs returns file extensions this converter can process
	// e.g., []string{".mov", ".avi", ".mkv"}
	SupportedInputs() []string

	// OutputExtension returns the output file extension
	// e.g., ".mp4"
	OutputExtension() string

	// Validate checks if the input file can be processed
	Validate(input string) error

	// Convert processes a single file with the given options
	Convert(input string, opts Options) (*Result, error)

	// ConvertBatch processes multiple files with the given options
	// Returns a slice of results and any fatal error
	ConvertBatch(inputs []string, opts Options) ([]*Result, error)
}

// ConverterWithSetup extends Converter with setup/teardown hooks
type ConverterWithSetup interface {
	Converter

	// Setup is called before any conversions begin
	Setup() error

	// Teardown is called after all conversions complete
	Teardown() error
}
