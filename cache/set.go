package cache

type set struct {
	dataCount uint32
	blocks    []*block
}

func (s *set) insert(ref, tag uint32, index int) {
	block := s.blocks[index]
	block.tag = tag
	block.data = int32(ref) // fake the data using the memory address as value
	block.validity = true
}
