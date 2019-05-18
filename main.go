package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joaofnds/cache-sim/cli"
	"github.com/joaofnds/cache-sim/file"
)

func main() {
	op := cli.Operation()
	switch op {
	case cli.GenerateData:
		err := genInputFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
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
	f, err := os.Create("input.dat")
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
	c, ref, err := cli.ParseNormalExecArgs(os.Args)
	if err != nil {
		switch err {
		case cli.ErrBadCacheFormat:
			fallthrough
		case cli.ErrBadArgNum:
			fmt.Println("Usage:")
			fmt.Println("\tcache_sim <nsets>:<bsize>:<assoc> input_file")
		}
		return err
	}

	b, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to parse cache to json: %v", err)
	}

	fmt.Printf("references: %v\ncache: %v\n", ref, string(b))

	return nil
}
