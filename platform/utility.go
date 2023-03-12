package platform

import (
	"errors"
	"fmt"
	"log"
	"net"
	"syscall"
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

func getDataLengthInBytes(length int32) []byte {
	length_byte_array := make([]byte, 4)

	length_byte_array[3] = byte(length & 255)
	length_byte_array[2] = byte((length >> 8) & 255)
	length_byte_array[1] = byte((length >> 16) & 255)
	length_byte_array[0] = byte((length >> 24) & 255)

	return length_byte_array[:]
}

func getLengthFromBytes(length_byte_array []byte) int {
	length := 0

	length |= int(length_byte_array[3])
	length |= int(length_byte_array[2]) << 8
	length |= int(length_byte_array[1]) << 16
	length |= int(length_byte_array[0]) << 24

	return length
}

func sendDataToNode(node *NetworkNode, data []byte, net_platform *NetworkPlatform) error {
	// connect to a network

	sending_addr, err := net.ResolveTCPAddr("tcp", node.GetAddressString())
	if err != nil {
		// TODO: Handle error
		log.Panic(err)
		return err
	}
	var conn *net.TCPConn
	if node == nil {
		conn, err = net.DialTCP("tcp", nil, sending_addr)
	} else if node.conn == nil {
		node.conn, err = net.DialTCP("tcp", nil, sending_addr)
		conn = node.conn
	} else {
		// log.Println("connection already exist")
		conn = node.conn
	}
	// conn, err = net.DialTCP("tcp", nil, sending_addr)
	if err != nil {
		fmt.Printf("Connection Failed, for node %s\n", node.Name)
		net_platform.AddToPreviousNodes(node.GetAddressString())
		net_platform.RemoveNodeWithAddress(node.GetAddressString())
		if net_platform.node_failure_event_handler != nil {
			net_platform.node_failure_event_handler(net_platform, node.GetAddressString())
		}
		return err
	}

	data = append(getDataLengthInBytes(int32(len(data))), data...)
	_, err = conn.Write(data)
	if err != nil {
		if errors.Is(err, syscall.EPIPE) {
			node.conn = nil
			return sendDataToNode(node, data, net_platform)
		}
		fmt.Printf("Sending data failed, error: %s\n", err.Error())
		return err
	}
	data = nil

	return nil
}

func sendDataToAddress(addr string, data []byte, net_platform *NetworkPlatform) error {

	// This is a blocking call make it non blocking
	sending_addr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		// TODO: Handle error
		return err
	}
	node_index := net_platform.get_node_from_string(addr)
	if node_index != -1 {
		return sendDataToNode(&net_platform.Connected_nodes[node_index], data, net_platform)
	}
	conn, err := net.DialTCP("tcp", nil, sending_addr)
	if err != nil {
		fmt.Printf("Connection Failed, for node with address: %s\nError: %s", addr, err)
		net_platform.AddToPreviousNodes(addr)
		net_platform.RemoveNodeWithAddress(addr)
		if net_platform.node_failure_event_handler != nil {
			net_platform.node_failure_event_handler(net_platform, addr)
		}
		//TODO: handle node failure
		return err
	}
	defer conn.Close()

	data = append(getDataLengthInBytes(int32(len(data))), data...)
	_, err = conn.Write(data) // write into connection i.e send data

	data = nil
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

		if addr_string[:3] == "192" || addr_string[:2] == "10" {
			return addr_string[:position]
		}
	}
	fmt.Print("Address not found")
	// return localhost if other is not found
	return "127.0.0.1"
}
