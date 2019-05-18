package cache

const (
	// ReplRandom is the random replacement strategy identifier
	ReplRandom = iota
)

// Cache is the struct that defines represents a cache. It is configured
// with a number of sets, block size and associativityhttps://github.com/zeromq/libzmq.
type Cache struct {
	Sets                int `json:"nsets"`
	BlockSize           int `json:"bsize"`
	Associativity       int `json:"assoc"`
	ReplacementStrategy int `json:"repl"`
}

// BuildCache builds a new cache with the given sets, block size, associativity
func BuildCache(sets, blockSize, associativity int) *Cache {
	newCache := Cache{}

	newCache.Sets = sets
	newCache.BlockSize = blockSize
	newCache.Associativity = associativity
	newCache.ReplacementStrategy = ReplRandom

	return &newCache
}
