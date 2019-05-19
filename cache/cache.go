package cache

import (
	"math"
)

const (
	// replRandom is the random replacement policy identifier
	replRandom = iota

	// ones an uint with the binary value full of ones, used to apply the index mask
	ones = ^uint32(0)
	// addressSize represents the bit length of memory addresses
	addressSize = uint32(32)
)

type block struct {
	validity bool
	tag      uint32
	data     int32
}

// Cache is the struct that defines represents a cache. It is configured
// with a number of sets, block size and associativity.
type Cache struct {
	cacheSize         uint32
	numberOfSets      uint32
	replacementPolicy uint8

	setSize uint32
	tagSize uint32

	indexMask uint32

	blocks []*block
}

// BuildCache builds a new cache with the given number of blocks,
// block size and associativity
func BuildCache(cacheSize, numberOfSets uint32) *Cache {
	// blocksPerSet := cacheSize / numberOfSets

	// directly mapped
	blocks := make([]*block, numberOfSets)
	for i := range blocks {
		blocks[i] = &block{}
	}

	setSize := uint32(math.Log2(float64(numberOfSets)))
	tagSize := addressSize - setSize

	indexMask := ones ^ (ones << setSize)

	return &Cache{
		cacheSize:         cacheSize,
		numberOfSets:      numberOfSets,
		replacementPolicy: replRandom,
		setSize:           setSize,
		tagSize:           tagSize,
		indexMask:         indexMask,
		blocks:            blocks,
	}
}

func (c *Cache) refSet(memoryReference uint32) uint32 {
	return (memoryReference & c.indexMask)
}

func (c *Cache) refTag(memoryReference uint32) uint32 {
	return (memoryReference >> c.setSize)
}

// Get retrieves data from the cache and inform if it was a hit or a miss
func (c *Cache) Get(ref uint32) (int32, bool) {
	index := c.refSet(ref)
	tag := c.refTag(ref)
	block := c.blocks[index]

	if block.tag != tag || !block.validity {
		// handle miss

		block.tag = tag
		block.data = int32(ref)
		block.validity = true

		return block.data, false
	}

	return block.data, true
}
