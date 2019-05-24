package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/joaofnds/cache-sim/cache"
	"github.com/joaofnds/cache-sim/file"
)

var (
	cliConfigRegexp = regexp.MustCompile(`^(\d+):(\d+):(\d)$`)
)

const (
	// GenerateData ...
	GenerateData = iota
	// NormalExecution ...
	NormalExecution
)

var (
	// ErrBadArgNum is the error returned when the number of arguments provided to the cli is wrong
	ErrBadArgNum = fmt.Errorf("wrong number of arguments. Must be 3")
	// ErrBadCacheFormat is the error return when the cache format provided is wrong
	ErrBadCacheFormat = fmt.Errorf("bad cache format. format is <nsets>:<bsize>:<assoc>")
)

// PrintUsage prints program usage to the provided io.Writer
func PrintUsage(w io.Writer) (int, error) {
	return fmt.Fprint(w, "Usage:\n\tcache_sim <nsets>:<bsize>:<assoc> input_file\n")
}

// Operation returns the operator to be performed by the program
func Operation() int {
	generate := flag.Bool("generate", false, "generate input data and exit")
	flag.Parse()

	if *generate {
		return GenerateData
	}

	return NormalExecution
}

// ParseNormalExecArgs parses command line args for the normal program execution
func ParseNormalExecArgs(args []string) (*cache.Cache, []uint32, error) {
	var c *cache.Cache
	var references []uint32
	if len(args) != 3 {
		return c, references, ErrBadArgNum
	}

	sets, blockSize, assoc, err := parseCacheConfig(args[1])
	if err != nil {
		return c, references, err
	}

	// TODO: use associativity instead of 1
	c = cache.BuildCache(sets, blockSize, assoc)

	fileName := args[2]
	f, err := os.Open(fileName)
	if err != nil {
		return c, references, err
	}
	defer f.Close()

	references, err = file.ParseInputFile(f)
	if err != nil {
		return c, references, err
	}

	return c, references, nil
}

// parseCacheConfig parses the cache config provided via command line
func parseCacheConfig(s string) (sets, blockSize, assoc uint32, err error) {
	matched := cliConfigRegexp.FindAllStringSubmatch(s, -1)
	if matched == nil || len(matched[0]) != 4 {
		err = ErrBadCacheFormat
		return
	}

	for k, v := range map[int]*uint32{1: &sets, 2: &blockSize, 3: &assoc} {
		n, e := strconv.Atoi(matched[0][k])
		if e != nil {
			err = ErrBadCacheFormat
			return
		}
		*v = uint32(n)
	}
	return
}
