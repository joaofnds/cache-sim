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

type report struct {
	accesses,
	hits,
	compulsoryMisses,
	capacityMisses,
	conflictMisses int
}

func (r *report) totalMisses() int {
	return r.compulsoryMisses + r.capacityMisses + r.conflictMisses
}

func main() {
	rand.Seed(time.Now().Unix())

	op := cli.Operation()
	switch op {
	case cli.Help:
		cli.PrintUsage(os.Stdout)
	case cli.GenerateData:
		err := genInputFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("'%s' generated!\n", file.DefaultFileName)
	case cli.RunSimulation:
		err := runSimulation()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
	}
	os.Exit(0)
}

func genInputFile() error {
	f, err := os.Create(file.DefaultFileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	s, n := cli.GenInputDataConfig()

	fmt.Printf("generating input file with %d entries in [0,%d)\n", s, n)
	if err := file.GenInputFile(f, s, n); err != nil {
		return fmt.Errorf("failed to generate input file: %v", err)
	}

	return nil
}

func runSimulation() error {
	c, addresses, err := cli.ParseSimulationArgs(os.Args)
	if err != nil {
		if err == cli.ErrBadArgNum {
			cli.PrintSimulationUsage(os.Stdout)
		}

		return err
	}

	r := report{accesses: len(addresses)}
	for _, address := range addresses {
		_, result := c.Get(address)

		if result == cache.Hit {
			r.hits++
		} else {
			switch result {
			case cache.MissCompulsory:
				r.compulsoryMisses++
			case cache.MissCapacity:
				r.capacityMisses++
			case cache.MissConflict:
				r.conflictMisses++
			}
		}
	}

	printReport(r)

	return nil
}

func printReport(r report) {
	fmt.Printf("accesses: %d\n", r.accesses)
	fmt.Printf("hits: %d\n", r.hits)
	fmt.Printf("misses: %d\n", r.totalMisses())
	fmt.Printf("compulsoryMisses: %d\n", r.compulsoryMisses)
	fmt.Printf("capacityMisses: %d\n", r.capacityMisses)
	fmt.Printf("conflictMisses: %d\n", r.conflictMisses)
}
