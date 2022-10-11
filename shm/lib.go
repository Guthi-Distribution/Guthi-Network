package shm

import (
	"fmt"
	"syscall"
	"unsafe"
)

/*
	Currently not used, we will be using unix package
*/

const (
	IPC_CREAT = 01000

	// Remove identifier.
	IPC_RMID = 0
	// Set `ipc_perm` options.
	IPC_SET = 1
	// Get `ipc_perm' options.
	IPC_STAT = 2
)

func ShmGet(key int, size int, mode int) (int, error) {
	id, _, errno := syscall.Syscall(syscall.SYS_SHMGET, uintptr(key), uintptr(size), uintptr(mode))
	if int(id) == -1 {
		return -1, errno
	}
	return int(id), nil
}

func ShmAt(id int, shm_addr uintptr, flag int) (*ShmSegment, error) {
	addr, _, errno := syscall.Syscall(syscall.SYS_SHMAT, uintptr(id), 0, uintptr(flag))
	if int(addr) == -1 {
		return nil, errno
	}

	shm_segment := (*ShmSegment)(unsafe.Pointer(addr))
	return shm_segment, nil
}

func ShmCtl(id int, cmd int) error {
	id_rem, _, err := syscall.Syscall(syscall.SYS_SHMCTL, uintptr(id), uintptr(IPC_RMID), uintptr(0))
	if int(id_rem) == -1 {
		fmt.Printf("Shared memory removal error: %s\n", err.Error())
		return err
	}

	return nil
}
