package platform

import (
	"GuthiNetwork/core"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"time"
)

/*
Passing of runtime info such as memory usage, CPU usage between nodes
*/

type CpuInformation struct {
	AddrFrom   string
	TotalUsage float32
	CpuStatus  core.ProcessorStatus
}

type MemoryInformation struct {
	AddrFrom     string
	MemoryStatus core.MemoryStatus
}

func SendGetCpuInfomation(addr string, net_platform *NetworkPlatform) error {
	payload := GetInformation{
		AddrFrom: net_platform.GetNodeAddress(),
	}
	data := GobEncode(payload)
	data = append(CommandStringToBytes("get_cpu_info"), data...)
	return sendDataToAddress(addr, data, net_platform)
}

func SendGetMemoryInfomation(addr string, net_platform *NetworkPlatform) error {
	payload := GetInformation{
		AddrFrom: net_platform.GetNodeAddress(),
	}
	data := GobEncode(payload)
	data = append(CommandStringToBytes("get_mem_info"), data...)
	return sendDataToAddress(addr, data, net_platform)
}

func HandleGetCpuInformation(request []byte, net_platform *NetworkPlatform) error {
	var payload GetInformation
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}
	// if the receiving address is the self address, then it is send
	send_payload := CpuInformation{
		net_platform.GetNodeAddress(),
		float32(core.GetCPUAllUsage()),
		core.GetProcessorInfo(),
	}

	return sendDataToAddress(payload.AddrFrom, append(CommandStringToBytes("cpuinfo"), GobEncode(send_payload)...), net_platform)
}

func HandleGetMemoryInformation(request []byte, net_platform *NetworkPlatform) error {
	var payload GetInformation
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}
	// if the receiving address is the self address, then it is send
	send_payload := MemoryInformation{
		net_platform.GetNodeAddress(),
		core.GetSysMemoryInfo(),
	}

	return sendDataToAddress(payload.AddrFrom, append(CommandStringToBytes("meminfo"), GobEncode(send_payload)...), net_platform)
}

func HandleReceiveMemoryInformation(request []byte, net_platfom *NetworkPlatform) error {
	var payload MemoryInformation
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}

	node_index := net_platfom.get_node_from_string(payload.AddrFrom)
	if node_index < 0 {
		//TODO: handle if the node information is not available, it can happen if the node is removed
		return nil
	}
	net_platfom.Connection_caches[node_index].Memory_info = payload.MemoryStatus

	return nil
}

func HandleReceiveCpuInformation(request []byte, net_platfom *NetworkPlatform) error {
	var payload CpuInformation
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}

	node_index := net_platfom.get_node_from_string(payload.AddrFrom)
	if node_index < 0 {
		// handle if the node information is not available, it can happen if the node is removed
		return nil
	}
	net_platfom.Connection_caches[node_index].Cpu_info = payload.CpuStatus
	net_platfom.Connection_caches[node_index].Cpu_info.Usage = payload.TotalUsage
	return nil
}

func RequestInfomation(net_platform *NetworkPlatform) {
	for true {
		time.Sleep(time.Second * 10)
		for _, node := range net_platform.Connected_nodes {
			SendGetCpuInfomation(node.GetAddressString(), net_platform)
			SendGetMemoryInfomation(node.GetAddressString(), net_platform)
		}
	}
}
