package shm

type ShmSegment struct {
	count uint16
	buff  []byte
}

func (memory *SharedMemory) WriteSharedMemory(data []byte) {
	length := uint16(len(data))
	memory.shm_segment.count = uint16(len(data))

	// TODO: Fix encoding
	// currently only memory alignment is used
	// the good old C way
	length_buff := []byte{
		byte(length & 0xff), byte((length >> 8) & 0xff),
	}
	data = append(length_buff, data...)
	copy(memory.shm_segment.buff, data)
}

func (memory *SharedMemory) ReadSharedMemory() []byte {
	memory.shm_segment.count = uint16(memory.shm_segment.buff[0]) | uint16(memory.shm_segment.buff[1])<<8

	return memory.shm_segment.buff[2:]
}
