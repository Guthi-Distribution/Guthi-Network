package core

import (
	"GuthiNetwork/shm"
	"unsafe"
)

//TODO: Need better name
/*
	Initializes the the basic structure needed for core communication
*/

/*
#cgo CXXFLAGS: "-std=c++20 -I../Guthi-Core/src/"
#cgo LDFLAGS: -L../Guthi-Core/lib -lGuthiCore -lpdh -lstdc++
#include "../Guthi-Core/src/core/c_api.h"
*/
import "C"

/*
Initializes the the basic structure needed for core communication
*/
func Initialize() error {
	initializeFileystem()
	var err error
	shared_memory, err = shm.CreateSharedMemory()
	if err != nil {
		return err
	}
	semaphore, err = shm.CreateSemaphore()
	if err != nil {
		return err
	}

	return nil
}

// Filesystem function
type FilesystemCore struct {
	Fs   string
	Size uint32
}

var filesystem FilesystemCore

func initializeFileystem() {
	data := C.GetLocalFileMetadata(unsafe.Pointer(&filesystem.Size))
	filesystem.Fs = C.GoString((*C.char)(*(*unsafe.Pointer)(data)))
}

func GetFileSystem() FilesystemCore {
	return filesystem
}

func SetFileSystem(fs FilesystemCore) {
	filesystem = fs
	shared_memory.WriteSharedMemory([]byte(filesystem.Fs), MESSSAGE_FILESYSTEM)
}

// Runtime info structure
// ------------------CPU----------------------
type ProcessorInfo struct {
	Processor_number uint32  `json:"processor_number"`
	Current_mhz      uint32  `json:"current_mhz"`
	Total_mhz        uint32  `json:"total_mhz"`
	Usage            float32 `json:"usage"`
}

type ProcessorStatus struct {
	Processor_count uint32          `json:"count"`
	Processors      []ProcessorInfo `json:"processors"`
}

func GetProcessorInfo() ProcessorStatus {
	info := C.GetSysProcessorInfo()
	status := ProcessorStatus{
		uint32(info.processor_count),
		[]ProcessorInfo{},
	}
	for i := 0; i < int(status.Processor_count); i++ {
		status.Processors = append(status.Processors, ProcessorInfo{
			uint32(info.processors[i].processor_number),
			uint32(info.processors[i].current_mhz),
			uint32(info.processors[i].total_mhz),
			float32(C.GetCurrentAllCPUUsage()),
		})
	}

	return status
}

/*
Returns total CPU usage in percentage
*/
func GetCPUAllUsage() float64 {
	usage := C.GetCurrentAllCPUUsage()
	return float64(usage)
}

// ------------------Memory--------------
type MemoryStatus struct {
	Installed_physical_ram uint64 `json:"installed"`
	Available_ram          uint64 `json:"Available"`
	Memory_load            uint64 `json:"Memory_Load"`
	// Information about virtual memory is not required here
}

func GetSysMemoryInfo() MemoryStatus {
	info := C.GetSysMemoryInfo()
	memory_status := MemoryStatus{
		uint64(info.installed_physical_ram),
		uint64(info.available_ram),
		uint64(info.memory_load),
	}

	return memory_status
}
