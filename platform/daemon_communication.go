package platform

import (
	"GuthiNetwork/daemon"
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
)

type DaemonMessage struct {
	AddrFrom string
	Message  daemon.MsgFormat
}

type FilesystemContents struct {
	AddrFrom string
	FileName string
	Contents []byte
}

func sendDaemonMessagesToNodes(message daemon.MsgFormat) error {
	payload := DaemonMessage{
		network_platform.GetNodeAddress(),
		message,
	}

	data := append(CommandStringToBytes("daemon_msg"), GobEncode(payload)...)

	for i := range network_platform.Connected_nodes {
		err := sendDataToNode(&network_platform.Connected_nodes[i], data, network_platform)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleDaemonMessageFromNodes(request []byte) {
	var payload DaemonMessage
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)

	daemon.SendFormattedMessage(*network_platform.daemon_handle, payload.Message)
}

func ListenForDaemonMessage() {
	for true {
		if network_platform.daemon_handle == nil {
			return
		}
		message := <-network_platform.daemon_handle.FormatChannel
		sendDaemonMessagesToNodes(message)

		fmt.Println("Message contents : ", string(message.Msg_content[:message.Msg_len]))
		switch message.Msg_type {
		case daemon.TrackedFileChanged:
			{
				fmt.Println("Tracked file changed : ", string(message.Msg_content[:message.Msg_len]))
				contents, err := ioutil.ReadFile(string(message.Msg_content[:message.Msg_len]))
				if err != nil {
					log.Println(err)
					continue
				}

				payload := FilesystemContents{
					network_platform.GetNodeAddress(),
					string(message.Msg_content[:message.Msg_len]),
					contents,
				}
				data := append(CommandStringToBytes("file_content"), GobEncode(payload)...)
				for i := range network_platform.Connected_nodes {
					sendDataToNode(&network_platform.Connected_nodes[i], data, network_platform)
				}
			}
		case daemon.TrackThisFile:
			{
				// Give absoute path
				// Not to be handled

			}
		case daemon.EchoMessage:
			{
				fmt.Println("Echo message received :-> ", string(message.Msg_content))
			}
		default:
			{
			}
		}
	}
}

func handleFileContents(request []byte) {
	var payload FilesystemContents
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)

	if network_platform.filesyste_merge != nil {
		network_platform.filesyste_merge(payload.Contents, payload.FileName)
	} else {
		err := ioutil.WriteFile(payload.FileName, payload.Contents, 0644)
		if err != nil {
			log.Printf("File writing error: %s\n", err)
		}
	}
}
