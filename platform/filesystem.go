package platform

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/Guthi/guthi_network/core"
)

func SendGetFileSystem(addr string, net_platform *NetworkPlatform) error {
	payload := GetInformation{
		AddrFrom: net_platform.GetNodeAddress(),
	}
	data := GobEncode(payload)
	data = append(CommandStringToBytes("get_fs"), data...)
	return sendDataToAddress(addr, data, net_platform)
}

func HandleGetFileSystem(request []byte, net_platform *NetworkPlatform) error {
	var payload GetInformation
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}
	// if the receiving address is the self address, then it is send
	send_payload := core.GetFileSystem()

	return sendDataToAddress(payload.AddrFrom, append(CommandStringToBytes("filesystem"), GobEncode(send_payload)...), net_platform)
}

func HandleReceiveFileSystem(request []byte, net_platfom *NetworkPlatform) error {
	var payload core.FilesystemCore
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}
	core.SetFileSystem(payload)

	return nil
}

/*
For now regularly file system is requested,
later it is sent when an event occurs
*/
func CommunicateFileSystem(net_platform *NetworkPlatform) {
	// for true {
	// 	time.Sleep(time.Second * 10)
	// 	for _,  := range net_platform.Connected_nodes {
	// 		// SendGetFileSystem(node.GetAddressString(), net_platform)
	// 	}
	// }
}
