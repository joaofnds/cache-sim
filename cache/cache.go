package cache

import "math/rand"

const (
	// ReplRandom is the random replacement strategy identifier
	ReplRandom = iota
)

// Entry defines a cache entry
type Entry struct {
	Validity bool   `json:"validity"`
	Ref      uint32 `json:"ref"`
	Data     int32  `json:"data"`
}

// Cache is the struct that defines represents a cache. It is configured
// with a number of sets, block size and associativity.
type Cache struct {
	Sets                int      `json:"nsets"`
	BlockSize           int      `json:"bsize"`
	Associativity       int      `json:"assoc"`
	ReplacementStrategy int      `json:"repl"`
	Entries             []*Entry `json:"entries"`
}

// Set sets data to the cache
func (c *Cache) Set(ref uint32, data int32) {
	i := c.refIndex(ref)
	entry := c.Entries[i]
	entry.Validity = true
	entry.Ref = ref
	entry.Data = data
}

// Get retrieves data from the cache
func (c *Cache) Get(ref uint32) (int32, bool) {
	i := c.refIndex(ref)
	entry := c.Entries[i]

	if entry.Ref != ref {
		newData := rand.Int31()
		c.Set(ref, newData)
		return newData, false
	}

	return entry.Data, true
}

func (c *Cache) refIndex(ref uint32) int {
	return int(ref) % c.Sets
}

// BuildCache builds a new cache with the given sets, block size, associativity
func BuildCache(sets, blockSize, associativity int) *Cache {
	entries := make([]*Entry, sets)

	for i := range entries {
		entries[i] = &Entry{}
	}

	return &Cache{
		Sets:                sets,
		BlockSize:           blockSize,
		Associativity:       associativity,
		ReplacementStrategy: ReplRandom,
		Entries:             entries,
	}
}
