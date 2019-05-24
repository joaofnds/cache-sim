package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/joaofnds/cache-sim/cache"
	"github.com/joaofnds/cache-sim/cli"
	"github.com/joaofnds/cache-sim/file"
)

func main() {
	rand.Seed(time.Now().Unix())

	op := cli.Operation()
	switch op {
	case cli.GenerateData:
		err := genInputFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("'%s' generated!\n", file.DefaultFileName)
		os.Exit(0)
	case cli.NormalExecution:
		err := normalExec()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	}
}

func genInputFile() error {
	f, err := os.Create(file.DefaultFileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	if err := file.GenInputFile(f); err != nil {
		return fmt.Errorf("failed to generate input file: %v", err)
	}

	return nil
}

func normalExec() error {
	c, refs, err := cli.ParseNormalExecArgs(os.Args)
	if err != nil {
		if err == cli.ErrBadArgNum {
			cli.PrintUsage(os.Stdout)
		}
		return err
	}

	for i := 0; i < 2; i++ {
		var hits, misses int
		for _, ref := range refs {
			if _, result := c.Get(ref); result == cache.Hit {
				hits++
			} else {
				misses++
			}
		}
		fmt.Printf("run %d:\n\thits: %d\n\tmisses: %d\n", i, hits, misses)
	}

	return nil
}
