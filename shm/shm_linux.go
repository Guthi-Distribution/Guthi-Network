package shm

import (
	"errors"
	"fmt"
	"syscall"

	"golang.org/x/sys/unix"
)

type SharedMemory struct {
	Id          int
	shm_segment ShmSegment
	key         int
	desc        *unix.SysvShmDesc
}

func CreateSharedMemory() (*SharedMemory, error) {
	var shm_memory SharedMemory
	shm_memory.shm_segment = ShmSegment{}
	shm_memory.desc = &unix.SysvShmDesc{}

	shm_memory.key = 69
	id, err := unix.SysvShmGet(shm_memory.key, 4098, (syscall.S_IRUSR | syscall.S_IWUSR | syscall.S_IRGRP | syscall.S_IWGRP))
	shm_memory.Id = id
	if err != nil {
		// a way for a specific error rather than just reading the error, and ofc go does not provide data of error variable
		if err.Error() == "no such file or directory" {
			return nil, errors.New("Shared Memory does not exist, Run core executable first")
		}
		fmt.Printf("Memory creation error: %s\n", err)
		return nil, err
	}

	buff, err := unix.SysvShmAttach(id, uintptr(0), 0)
	if err != nil {
		fmt.Printf("Memory Attachment error: %s\n", err)
		return nil, err
	}
	shm_memory.shm_segment.buff = buff

	return &shm_memory, nil
}

func (memory *SharedMemory) RemoveSharedMemory() error {
	err := unix.SysvShmDetach(memory.shm_segment.buff)
	if err != nil {
		fmt.Printf("Error detaching memory: %s", err)
		return err
	}
	_, err = unix.SysvShmCtl(memory.Id, unix.IPC_RMID, nil)
	if err != nil {
		fmt.Printf("Error Removing memory: %s", err)
		return err
	}

	return nil
}
