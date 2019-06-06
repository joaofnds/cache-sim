package cli

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/joaofnds/cache-sim/cache"
	"github.com/joaofnds/cache-sim/file"
)

var (
	// ErrBadArgNum is the error returned when the number of arguments provided to the cli is wrong
	ErrBadArgNum = fmt.Errorf("wrong number of arguments. Must be 3")
	// ErrBadCacheFormat is the error return when the cache format provided is wrong
	ErrBadCacheFormat = fmt.Errorf("bad cache format. format is <nsets>:<bsize>:<assoc>")
	// ErrBadBlockSize is the error returned when the specified block size is wrong
	ErrBadBlockSize = fmt.Errorf("bad block size. must be { 2^x | x >= 5 }")

	cliConfigRegexp = regexp.MustCompile(`^(\d+):(\d+):(\d+)$`)

	cacheFormatUsage = "cache_sim <nsets>:<bsize>:<assoc> input_file"
)

// PrintSimulationUsage prints program usage to the provided io.Writer
func PrintSimulationUsage(w io.Writer) (int, error) {
	return fmt.Fprintf(w, "Usage:\n\t%s\n", cacheFormatUsage)
}

// ParseSimulationArgs parses command line args to run the simulation
func ParseSimulationArgs(args []string) (*cache.Cache, []uint32, error) {
	var c *cache.Cache
	var addresses []uint32
	if len(args) != 3 {
		return c, addresses, ErrBadArgNum
	}

	sets, blockSize, assoc, err := parseCacheConfig(args[1])
	if err != nil {
		return c, addresses, err
	}

	c = cache.BuildCache(sets, blockSize, assoc)

	fileName := args[2]
	f, err := os.Open(fileName)
	if err != nil {
		return c, addresses, err
	}
	defer f.Close()

	addresses, err = file.ParseInputFile(f)
	if err != nil {
		return c, addresses, err
	}

	return c, addresses, nil
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

	if !cache.ValidBlockSize(blockSize) {
		err = ErrBadBlockSize
		return
	}

	return
}
