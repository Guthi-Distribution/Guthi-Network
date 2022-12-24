package shm

type ShmSegment struct {
	count uint16
	buff  []byte
}

type MESSAGE_TYPE uint8

func (memory *SharedMemory) WriteSharedMemory(data []byte, message MESSAGE_TYPE) {
	length := uint16(len(data))
	memory.shm_segment.count = uint16(len(data))

	message_type := byte(message + '0')
	length_buff := []byte{
		byte(length & 0xff), byte((length >> 8) & 0xff), '1', message_type,
	}
	data = append(length_buff, data...)
	copy(memory.shm_segment.buff, data)
}

func (memory *SharedMemory) ReadSharedMemory() []byte {
	memory.shm_segment.count = uint16(memory.shm_segment.buff[0]) | uint16(memory.shm_segment.buff[1])<<8

	return memory.shm_segment.buff[2:]
}
