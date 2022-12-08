package platform

import (
	"bytes"
	"errors"
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
	fmt.Printf("Hello there")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("Connection Failed, for node with address: %s\nError: %s", addr, err)
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

func getForwardSlashPosition(value string) int {
	for i, c := range value {
		if c == '/' {
			return i
		}
	}

	return -1
}

func GetNodeAddress() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		log.Panic(err.Error())
	}

	for _, addr := range addresses {
		addr_string := addr.String()
		position := getForwardSlashPosition(addr_string)

		if addr_string[:3] == "192" {
			return addr_string[:position]
		}
	}

	err = errors.New("Address not found")
	log.Panic(err)
	return ""
}
