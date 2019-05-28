package cache

type set struct {
	dataCount uint32
	blocks    []*block
}

func (s *set) insert(index int, tag uint32, data []int32) {
	block := s.blocks[index]
	block.tag = tag
	block.data = data // fake the data using the memory address as value
	block.validity = true
}
