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
	processor_number uint32
	current_mhz      uint32
	total_mhz        uint32
}

type ProcessorStatus struct {
	processor_count uint32
	processors      []ProcessorInfo
}

func GetProcessorInfo() ProcessorStatus {
	info := C.GetSysProcessorInfo()
	status := ProcessorStatus{
		uint32(info.processor_count),
		[]ProcessorInfo{},
	}
	for i := 0; i < int(status.processor_count); i++ {
		info := ProcessorInfo{
			uint32(info.processors[i].processor_number),
			uint32(info.processors[i].current_mhz),
			uint32(info.processors[i].total_mhz),
		}
		status.processors = append(status.processors, info)
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
	installed_physical_ram uint64
	available_ram          uint64
	memory_load            uint64
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
