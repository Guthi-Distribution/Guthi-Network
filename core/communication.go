package core

import "GuthiNetwork/shm"

/*
Shared memory format:
  - First byte contain the runtimen type: 0 for Go and 1 for C++, if the value is 1 then it can be interpreted as C++ daemon wrote in it and Go has to read it
  - Second byte byte contain message type (event, or data)
  - For later data we need format for each message type
*/
const (
	RUNTIME_TIME = 0 // FOR GO

	// Message Type
	MESSAGE_EVENT = 0

	// EVENT TYPES
	FILE_CHANGED_EVENT = 0
)

var semaphore *shm.Semaphore
var shared_memory *shm.SharedMemory
