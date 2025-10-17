package converter

import (
	"context"
	"time"
)

// Options contains common conversion options
type Options struct {
	// Input/Output
	Input     string
	Output    string
	OutputDir string

	// Processing
	Workers       int
	SkipExisting  bool
	DryRun        bool
	Verbose       bool
	FlatStructure bool

	// Progress
	ShowProgress bool
	Context      context.Context
}

// Result represents the outcome of a conversion
type Result struct {
	Input       string
	Output      string
	Success     bool
	Error       error
	Duration    time.Duration
	InputSize   int64
	OutputSize  int64
	Skipped     bool
	SkipReason  string
}

// Stats tracks conversion statistics
type Stats struct {
	Total     int
	Success   int
	Failed    int
	Skipped   int
	StartTime time.Time
	EndTime   time.Time
}

// Progress represents ongoing conversion progress
type Progress struct {
	Current     int
	Total       int
	CurrentFile string
	Percentage  float64
}
