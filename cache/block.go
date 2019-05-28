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
	log := math.Log2(float64(n))
	return log >= 5 && math.Mod(log, 1) == 0
}

func (b *block) get(ref uint32) int32 {
	dataLen := len(b.data)
	shiftSize := uint32(math.Log2(float64(dataLen)))
	index := int(ref&(^shiftSize)) % dataLen // first bit

	return b.data[index]
}
