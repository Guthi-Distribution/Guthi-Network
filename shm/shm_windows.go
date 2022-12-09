package shm

import (
	"errors"
	"fmt"
	"log"
	"unsafe"

	"golang.org/x/sys/windows"
)

/*
#include <Windows.h>
#include <strsafe.h>
char* get_error_message(const char* message) {
    char* error_message_buff;
    DWORD dw = GetLastError();

    FormatMessage(
        FORMAT_MESSAGE_ALLOCATE_BUFFER |
        FORMAT_MESSAGE_FROM_SYSTEM |
        FORMAT_MESSAGE_IGNORE_INSERTS,
        message,
        dw,
        MAKELANGID(LANG_NEUTRAL, SUBLANG_DEFAULT),
        (LPTSTR)&error_message_buff,
        0, NULL);
    char* error = (char*)malloc(1024);
    sprintf_s(error, 1024, "%s:%s", message, error_message_buff);
    LocalFree(error_message_buff);
    return error;
}

typedef struct ShmInfo {
    HANDLE hnd;
    char* err; // to make this stuff "goish"
}ShmInfo;

ShmInfo* OpenSharedMemory(char *name) {
    ShmInfo *info = (ShmInfo *)malloc(sizeof(ShmInfo));
    info->hnd = OpenFileMapping (
        FILE_MAP_ALL_ACCESS,   // read/write access
        FALSE,                 // do not inherit the name
        name
    );
    info->err = NULL;
    if (info->hnd == NULL) {
        info->err = get_error_message("Opening of shared memory failed");
        return info;
    }

    return info;
}
*/
import "C"

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
	// name, err := windows.UTF16PtrFromString("Guthi_Shared_memory")
	// if err != nil {
	// 	fmt.Printf("String to pointer conversion error: %s", err)
	// }
	// hnd, err := windows.CreateFileMapping(
	// 	windows.InvalidHandle,
	// 	nil,
	// 	windows.PAGE_READWRITE,
	// 	0,
	// 	4100,
	// 	name,
	// )
	// if err != nil {
	// 	log.Fatalf("CreateFileMapping error: %s", err.Error())
	// 	return nil, err
	// }
	name := "guthi_semaphore"
	c_mem := C.OpenSharedMemory(C.CString(name))
	err := C.GoString(c_mem.err)
	if c_mem == nil || c_mem.hnd == nil {
		return nil, errors.New(err)
	}
	hnd := windows.Handle(unsafe.Pointer(c_mem.hnd))
	addr, e := windows.MapViewOfFile(hnd, FILE_MAP_ALL_ACCESS, 0, 0, 4100)
	if e != nil {
		log.Fatalf("MapViewOfFile error: %s", e.Error())
		return nil, e
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
