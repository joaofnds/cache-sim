package main

import (
	"flag"
	"log"
	"os"

	"github.com/joaofnds/cache-sim/file"
)

func main() {
	generate := flag.Bool("generate", false, "generate input data and exit")
	flag.Parse()

	if *generate {
		genInputFile()
		return
	}
}

func genInputFile() {
	f, err := os.Create("input.dat")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	if err := file.GenInputFile(f); err != nil {
		log.Fatalf("failed to generate input file: %v\n", err)
	}
}
