package core

/*
#cgo CXXFLAGS: "-std=c++20 -I../Guthi-Core/src/"
#cgo LDFLAGS: -L../Guthi-Core/ -lGuthiCore -lstdc++
#include "../Guthi-Core/src/core/c_api.h"
*/
import "C"
import "fmt"

func Initialize() interface{} {
	GetSysMemoryInfo()
	fmt.Println(GetCPUAllUsage())
	fmt.Println(GetCPUAllUsage())
	return GetProcessorInfo()
}

// Runtime info structure

// ------------------CPU----------------------
type ProcessorInfo struct {
	Processor_number uint32
	Current_mhz      uint32
	Total_mhz        uint32
}

type ProcessorStatus struct {
	Processor_count uint32
	Processors      []ProcessorInfo
}

func GetProcessorInfo() ProcessorStatus {
	info := C.GetSysProcessorInfo()
	status := ProcessorStatus{
		uint32(info.processor_count),
		[]ProcessorInfo{},
	}
	for i := 0; i < int(status.Processor_count); i++ {
		info := ProcessorInfo{
			uint32(info.processors[i].processor_number),
			uint32(info.processors[i].current_mhz),
			uint32(info.processors[i].total_mhz),
		}
		status.Processors = append(status.Processors, info)
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
