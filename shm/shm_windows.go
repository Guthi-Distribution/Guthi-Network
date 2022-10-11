package shm

import (
	"fmt"
	"log"
	"unsafe"

	"golang.org/x/sys/windows"
)

type SharedMemory struct {
	hnd         windows.Handle
	shm_segment ShmSegment
}

const (
	length                       = 4100
	name                         = "Guthi_Shared_memory"
	SECTION_QUERY                = 0x0001
	SECTION_MAP_WRITE            = 0x0002
	SECTION_MAP_READ             = 0x0004
	SECTION_MAP_EXECUTE          = 0x0008
	SECTION_EXTEND_SIZE          = 0x0010
	SECTION_MAP_EXECUTE_EXPLICIT = 0x0020
	FILE_MAP_ALL_ACCESS          = windows.STANDARD_RIGHTS_REQUIRED |
		SECTION_QUERY |
		SECTION_MAP_WRITE |
		SECTION_MAP_READ |
		SECTION_MAP_EXECUTE |
		SECTION_EXTEND_SIZE
)

func CreateSharedMemory() (*SharedMemory, error) {
	memory := &SharedMemory{}
	name, err := windows.UTF16PtrFromString("Guthi_Shared_memory")
	if err != nil {
		fmt.Printf("String to pointer conversion error: %s", err)
	}
	hnd, err := windows.CreateFileMapping(
		windows.InvalidHandle,
		nil,
		windows.PAGE_READWRITE,
		0,
		4100,
		name,
	)
	if err != nil {
		log.Fatalf("CreateFileMapping error: %s", err.Error())
		return nil, err
	}

	addr, err := windows.MapViewOfFile(hnd, FILE_MAP_ALL_ACCESS, 0, 0, 4100)
	if err != nil {
		log.Fatalf("MapViewOfFile error: %s", err.Error())
		return nil, err
	}

	memory.hnd = hnd
	var b = struct {
		addr uintptr
		len  int
		cap  int
	}{addr, int(length), int(length)}

	memory.shm_segment.buff = *(*[]byte)(unsafe.Pointer(&b))
	return memory, nil
}

func (memory *SharedMemory) RemoveSharedMemory() error {
	err := windows.UnmapViewOfFile(uintptr(unsafe.Pointer(&memory.shm_segment.buff[0])))
	if err != nil {
		fmt.Printf("Unmap of memory error: %s\n", err)
		return err
	}
	err = windows.CloseHandle(memory.hnd)
	if err != nil {
		fmt.Printf("Closing of handler error: %s\n", err)
	}

	return nil
}
