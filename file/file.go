// Package file deals with cache-sim file IO
package file

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
)

const (
	// DefaultFileName is the name of the generated input file
	DefaultFileName = "input.dat"
)

// GenInputFile generates random input data
func GenInputFile(f io.Writer, size int, n int32) error {
	for i := 0; i < size; i++ {
		err := binary.Write(f, binary.BigEndian, rand.Int31n(n))
		if err != nil {
			return fmt.Errorf("failed to write integer to file: %f", err)
		}
	}

	return nil
}

// ParseInputFile parses an input file in the same format as the generated
// output file, a binary file containing consecutive 32bit integers
func ParseInputFile(r io.Reader) ([]uint32, error) {
	inputs := make([]uint32, 0)
	var n uint32
	for {
		err := binary.Read(r, binary.BigEndian, &n)
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("failed while reading bytes: %v", err)
		}

		inputs = append(inputs, n)
	}

	return inputs, nil
}
