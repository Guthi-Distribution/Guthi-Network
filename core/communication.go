package core

import (
	"GuthiNetwork/shm"
	"fmt"
)

/*
Shared memory format:
  - First byte contain the runtimen type: 0 for Go and 1 for C++, if the value is 1 then it can be interpreted as C++ daemon wrote in it and Go has to read it
  - Second byte byte contain message type (event, or data)
  - For later data we need format for each message type
*/
const (
	RUNTIME_TYPE = 1 // FOR GO process id

	// Message Type
	MESSAGE_EVENT       = 0
	MESSSAGE_FILESYSTEM = 1

	// EVENT TYPES
	FILE_CHANGED_EVENT = 0
)

var semaphore *shm.Semaphore
var shared_memory *shm.SharedMemory

func ReadSharedMemory() {
	data := shared_memory.ReadSharedMemory()
	// 0x30 is the acii value of
	if data[0] == 0x30 {
		data = data[1:]
		filesystem.Fs = string(data[:])
		fmt.Println(filesystem.Fs)
	}
}
