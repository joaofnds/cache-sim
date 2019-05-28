package cache

import (
	"math"
	"math/rand"
)

const (
	// replRandom is the random replacement policy identifier
	replRandom = iota

	// Hit indicates that the cache access was a hit
	Hit = iota
	// MissCompulsory indicates that the cache access had a compulsory miss
	MissCompulsory
	// MissCapacity indicates that the cache access had a capacity miss
	MissCapacity
	// MissConflict indicates that the cache access has a conflict miss
	MissConflict

	// ones is an uint with the binary value full of ones, used to apply the index mask
	ones = ^uint32(0)
	// addressSize represents the bit length of memory addresses
	addressSize = uint32(32)
)

// Cache is the struct that represents a cache. It is configured
// with a number of sets, block size and associativity.
type Cache struct {
	cacheSize         uint32
	numberOfSets      uint32
	replacementPolicy uint8

	// set info
	indexSize uint32
	tagSize   uint32
	indexMask uint32
	assoc     uint32

	sets []*set
}

// BuildCache builds a new cache given the number of blocks,
// block size and associativity
func BuildCache(numberOfSets, blockSize, assoc uint32) *Cache {
	indexSize := uint32(math.Log2(float64(assoc)))
	tagSize := addressSize - indexSize
	indexMask := ones ^ (ones << indexSize)

	sets := make([]*set, numberOfSets)
	for i := range sets {
		sets[i] = &set{
			dataCount: 0,
		}

		blocks := make([]*block, assoc)
		for i := range blocks {
			blocks[i] = &block{}
		}

		sets[i].blocks = blocks
	}

	return &Cache{
		cacheSize:         uint32(0),
		numberOfSets:      numberOfSets,
		replacementPolicy: replRandom,
		sets:              sets,

		tagSize:   tagSize,
		indexMask: indexMask,
		assoc:     assoc,
	}
}

// Get retrieves data from the cache and inform if it was a hit or a miss
func (c *Cache) Get(ref uint32) (int32, int) {
	setIndex := ref % c.numberOfSets

	return c.getOnSet(c.sets[setIndex], ref)
}

func accessResult(block *block, tag uint32) int {
	if block.tag == 0 && !block.validity && block.data == 0 {
		return MissCompulsory
	}

	if block.tag != tag {
		return MissConflict
	}

	return Hit
}

// refSet returns the set index of the given memoryReference
func (c *Cache) refSet(memoryReference uint32) uint32 {
	return (memoryReference & c.indexMask)
}

// refTag returns the tag of the given memoryReference
func (c *Cache) refTag(memoryReference uint32) uint32 {
	return (memoryReference >> c.indexSize)
}

// handleMiss "get the date" and insert it into the set
func (c *Cache) handleMiss(set *set, ref, tag uint32) {
	if set.dataCount < c.assoc {
		// set with available spaces
		for i, block := range set.blocks {
			if !block.validity {
				set.insert(ref, tag, i)
				break
			}
		}
		set.dataCount++
	} else {
		// set full
		i := rand.Intn(int(c.assoc))
		set.insert(ref, tag, i)
	}
}

// Get retrieves data from the cache and inform if it was a hit or a miss
func (c *Cache) getOnSet(set *set, ref uint32) (int32, int) {
	tag := c.refTag(ref)

	var block *block
	for _, block = range set.blocks {
		if block.tag == tag {
			break
		}
	}

	hitOrMiss := accessResult(block, tag)

	if hitOrMiss != Hit {
		c.handleMiss(set, ref, tag)
	}

	return block.data, hitOrMiss
}
