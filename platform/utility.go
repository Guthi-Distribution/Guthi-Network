package platform

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func CommandStringToBytes(cmd string) []byte {
	var bytes [COMMAND_LENGTH]byte

	for i, c := range cmd {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func BytesToCommandString(bytes []byte) string {
	var cmd []byte

	for _, b := range bytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}

	return fmt.Sprintf("%s", cmd)
}

func sendDataToNode(node *NetworkNode, data []byte, net_platform *NetworkPlatform) {
	// connect to a network
	conn, err := net.Dial(node.Socket.Network(), node.Socket.String())

	if err != nil {
		log.Printf("Connection Failed, for node %s\n", node.Name)
		net_platform.RemoveNode(*node)
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Printf("Sending data failed, error: %s\n", err.Error())
	}
}

func sendDataToAddress(addr string, data []byte, net_platform *NetworkPlatform) error {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Printf("Connection Failed, for node with address: %s\n", addr)
		net_platform.AddToPreviousNodes(addr)
		net_platform.RemoveNodeWithAddress(addr)
		//TODO: handle node failure
		return err
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data)) // write into connection i.e send data

	if err != nil {
		return err
	}

	return err
}
