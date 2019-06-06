package cache

type block struct {
	validity bool
	tag      uint32
	data     []byte
}
