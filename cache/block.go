package cache

import (
	"math"
)

type block struct {
	validity bool
	tag      uint32
	data     []int32
}

// ValidBlockSize validates the cache block size
func ValidBlockSize(n uint32) bool {
	bitSize := math.Log2(float64(n))
	return n >= addressSize && math.Mod(bitSize, 1) == 0
}
