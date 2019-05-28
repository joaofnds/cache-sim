package cli

import (
	"fmt"
	"io"
	"os"
)

const (
	// Noop represents no operation
	Noop = iota
	// Help command prints program usage
	Help
	// GenerateData ...
	GenerateData
	// RunSimulation ...
	RunSimulation
)

// Operation returns the operator to be performed by the program
func Operation() int {
	if len(os.Args) < 2 {
		return Help
	}

	switch os.Args[1] {
	case "-help":
		fallthrough
	case "--help":
		fallthrough
	case "help":
		return Help
	case "generate":
		parseGenInputDataConfig()
		return GenerateData
	default:
		return RunSimulation
	}
}

// PrintUsage prints program usage to the provided io.Writer
func PrintUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintf(w, "\t%s\n", cacheFormatUsage)
	fmt.Fprintln(w, "\tcache-sim generate -size int -range int")
}
