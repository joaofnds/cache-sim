package cli

import (
	"flag"
	"os"
)

var (
	genInputSize  int
	genInputRange int32
)

// ParseGenInputDataConfig returns the input file config
func GenInputDataConfig() (int, int32) {
	return genInputSize, genInputRange
}

func parseGenInputDataConfig() {
	genSet := flag.NewFlagSet("generate", flag.ExitOnError)
	inputSize := genSet.Int("size", 32, "size of the input file in entries")
	inputRange := genSet.Int("range", 16, "range of entries. [0,range)")

	genSet.Parse(os.Args[2:])
	genInputSize = *inputSize
	genInputRange = int32(*inputRange)
}
