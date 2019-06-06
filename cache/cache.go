// Package cache implements the actual cache-sim cache
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
	// wordSize represents the bit length of a word
	wordSize = uint32(8)
)

// Cache is the struct that represents a cache. It is configured
// with a number of sets, block size and associativity.
type Cache struct {
	// cacheSize ...
	cacheSize uint32
	// numberOfSets ...
	numberOfSets uint32
	// replacementPolicy ...
	replacementPolicy uint8

	// set info

	// indexSize number of bits needed to store the set index
	indexSize uint32
	// tagSize number of bits needed to store the tag
	tagSize uint32
	// indexMask mask used to get the set index a memory address
	indexMask uint32
	// assoc degree of associativity
	assoc uint32

	// block info

	// blockSize is the number of words stored in a blcok
	blockSize uint32
	// offsetMask is a mask used to get the index of a word in the block
	offsetMask uint32

	// sets are the actual sets in the cache
	sets []*set
}

// BuildCache builds a new cache given the number of blocks,
// block size and associativity
func BuildCache(numberOfSets, blockSize, assoc uint32) *Cache {
	indexSize := uint32(math.Log2(float64(assoc)))
	tagSize := addressSize - indexSize
	indexMask := ones ^ (ones << indexSize)
	offsetSize := uint32(math.Log2(float64(blockSize)))
	offsetMask := ^offsetSize

	sets := make([]*set, numberOfSets)
	for i := range sets {
		sets[i] = &set{
			dataCount: 0,
		}

		blocks := make([]*block, assoc)
		for i := range blocks {
			blocks[i] = &block{
				data: make([]byte, blockSize),
			}
		}

		sets[i].blocks = blocks
	}

	return &Cache{
		cacheSize:         uint32(0),
		numberOfSets:      numberOfSets,
		replacementPolicy: replRandom,

		tagSize:   tagSize,
		indexMask: indexMask,
		assoc:     assoc,

		blockSize:  blockSize,
		offsetMask: offsetMask,

		sets: sets,
	}
}

// Get retrieves data from the cache and inform if it was a hit or a miss
func (c *Cache) Get(ref uint32) (byte, int) {
	setIndex := (ref / c.blockSize) % c.numberOfSets

	return c.getOnSet(c.sets[setIndex], ref)
}

// accessResult checks if the access to the cache is a hit or a miss.
// It also informs what kind of miss was it.
func accessResult(block *block, tag uint32) int {
	if block.tag == 0 && !block.validity {
		return MissCompulsory
	}

	if block.tag != tag {
		return MissConflict
	}

	return Hit
}

// addressSet returns the set index of the given memory address
func (c *Cache) addressSet(address uint32) uint32 {
	return address & c.indexMask
}

// addressTag returns the tag of the given memory address
func (c *Cache) addressTag(address uint32) uint32 {
	blockAddress := address - (address % c.blockSize)
	return blockAddress >> c.indexSize
}

// handleMiss "get the date" and insert it into the set
func (c *Cache) handleMiss(set *set, ref, tag uint32) {
	if set.dataCount < c.assoc {
		// set with available spaces
		for i, block := range set.blocks {
			if !block.validity {
				data := c.retrieveFromLowerLevel(ref)
				set.insert(i, tag, data)
				break
			}
		}
		set.dataCount++
	} else {
		// set full
		i := rand.Intn(int(c.assoc))
		data := c.retrieveFromLowerLevel(ref)
		set.insert(i, tag, data)
	}
}

// retrieveFromLowerLevel emulates the action of retrieving the
// data from a lower memory in the hiearchy
func (c *Cache) retrieveFromLowerLevel(memoryAddress uint32) []byte {
	data := make([]byte, c.blockSize)

	blockStart := int32(memoryAddress - (memoryAddress % c.blockSize))
	blockEnd := blockStart + int32(c.blockSize)

	for address := blockStart; address < blockEnd; address++ {
		i := address % int32(c.blockSize)
		data[i] = byte(address) // fake the retrieval by storing the address itself as the data
	}

	return data
}

// Get retrieves data from the cache and inform if it was a hit or a miss
func (c *Cache) getOnSet(set *set, address uint32) (byte, int) {
	tag := c.addressTag(address)

	var block *block
	for _, block = range set.blocks {
		if block.tag == tag {
			break
		}
	}

	hitOrMiss := accessResult(block, tag)

	if hitOrMiss != Hit {
		c.handleMiss(set, address, tag)
	}

	data := c.getOnBlock(block, address)

	return data, hitOrMiss
}

// getOnBlock retrieves the data from a block
func (c *Cache) getOnBlock(b *block, address uint32) byte {
	index := int(address&c.offsetMask) % int(c.blockSize)

	return b.data[index]
}
