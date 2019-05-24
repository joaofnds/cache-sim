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

type block struct {
	validity bool
	tag      uint32
	data     int32
}

type set struct {
	indexSize uint32
	tagSize   uint32
	indexMask uint32
	assoc     uint32
	dataCount uint32
	blocks    []*block
}

func (s *set) refSet(memoryReference uint32) uint32 {
	return (memoryReference & s.indexMask)
}

func (s *set) refTag(memoryReference uint32) uint32 {
	return (memoryReference >> s.indexSize)
}

func (s *set) insert(ref, tag uint32, index int) {
	block := s.blocks[index]
	block.tag = tag
	block.data = int32(ref)
	block.validity = true
}

func (s *set) handleMiss(ref, tag uint32) {
	if s.dataCount < s.assoc {
		// set with available spaces
		for i, block := range s.blocks {
			if !block.validity {
				s.insert(ref, tag, i)
				break
			}
		}
		s.dataCount++
	} else {
		// set full
		i := rand.Intn(int(s.assoc))
		s.insert(ref, tag, i)
	}
}

// Get retrieves data from the cache and inform if it was a hit or a miss
func (s *set) get(ref uint32) (int32, int) {
	tag := s.refTag(ref)

	var block *block
	for _, block = range s.blocks {
		if block.tag == tag {
			break
		}
	}

	hitOrMiss := accessResult(block, tag)

	if hitOrMiss != Hit {
		s.handleMiss(ref, tag)
	}

	return block.data, hitOrMiss
}

// Cache is the struct that defines represents a cache. It is configured
// with a number of sets, block size and associativity.
type Cache struct {
	cacheSize         uint32
	numberOfSets      uint32
	replacementPolicy uint8

	sets []*set
}

// BuildCache builds a new cache with the given number of blocks,
// block size and associativity
func BuildCache(numberOfSets, blockSize, assoc uint32) *Cache {
	indexSize := uint32(math.Log2(float64(assoc)))
	tagSize := addressSize - indexSize
	indexMask := ones ^ (ones << indexSize)

	sets := make([]*set, numberOfSets)
	for i := range sets {
		sets[i] = &set{
			tagSize:   tagSize,
			indexMask: indexMask,
			assoc:     assoc,
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
	}
}

// Get retrieves data from the cache and inform if it was a hit or a miss
func (c *Cache) Get(ref uint32) (int32, int) {
	setIndex := ref % c.numberOfSets

	return c.sets[setIndex].get(ref)
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
