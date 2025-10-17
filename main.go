package main

import (
	"github.com/onedusk/sb/cmd"
	"github.com/onedusk/sb/cmd/formats"
	_ "github.com/onedusk/sb/internal/processors/mov_to_mp4" // Register MP4 converter
)

func main() {
	// Register format commands
	cmd.GetRootCmd().AddCommand(formats.MP4Cmd)

	// Execute CLI
	cmd.Execute()
}
