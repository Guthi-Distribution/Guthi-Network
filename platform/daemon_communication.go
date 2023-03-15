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

	// daemon.SendFormattedMessage(*network_platform.daemon_handle, payload.Message)
}

func ListenForDaemonMessage() {
	for true {
		if network_platform.daemon_handle == nil {
			return
		}
		message := <-network_platform.daemon_handle.FormatChannel
		// sendDaemonMessagesToNodes(message)

		fmt.Println("Message contents : ", string(message.Msg_content[:message.Msg_len]))
		switch message.Msg_type {
		case daemon.TrackedFileChanged:
			{
				fmt.Printf("Tracked file changed: %s Length: %d", string(message.Msg_content[:message.Msg_len]), message.Msg_len)
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
				log.Println(len(network_platform.Connected_nodes))
				for i := range network_platform.Connected_nodes {
					log.Printf("Sending File content: %s\n", network_platform.Connected_nodes[i].Name)
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
		log.Println("WARN: No callback function for handling file merge")
		err := ioutil.WriteFile(payload.FileName, payload.Contents, 0644)

		if err != nil {
			log.Printf("File writing error: %s\n", err)
		}
	}

	log.Printf("Tracking file:%s\n", payload.FileName)
	network_platform.TrackFile(payload.FileName)
}
