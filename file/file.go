// Package file deals with cache-sim file IO
package file

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"time"
)

// GenInputFile generates random input data
func GenInputFile(f io.Writer) error {

	rand.Seed(time.Now().Unix())

	bounds := [2]int32{10, 1000}
	for _, bound := range bounds {
		for i := 0; i < 100; i++ {
			err := binary.Write(f, binary.BigEndian, rand.Int31n(bound))
			if err != nil {
				return fmt.Errorf("failed to write integer to file: %f", err)
			}
		}
	}

	return nil
}

func parse(r io.Reader) ([]int32, error) {
	inputs := make([]int32, 0)
	var n int32
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
