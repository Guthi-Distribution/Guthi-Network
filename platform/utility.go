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
		fmt.Printf("Connection Failed, for node %s\n", node.Name)
		net_platform.RemoveNode(*node)
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	data = nil
	if err != nil {
		fmt.Printf("Sending data failed, error: %s\n", err.Error())
	}
}

func sendDataToAddress(addr string, data []byte, net_platform *NetworkPlatform) error {
	// This is a blocking call make it non blocking
	conn, err := net.Dial("tcp", addr)
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

	_, err = io.Copy(conn, bytes.NewReader(data)) // write into connection i.e send data

	if err != nil {
		return err
	}

	return err
}

// func sendDataToAddress(addr string, data []byte, net_platform *NetworkPlatform) error {
// 	// This is a blocking call make it non blocking
// 	// conn, err := net.Dial("tcp", addr)
// 	var err error
// 	index := net_platform.get_node_from_string(addr)
// 	var conn *net.TCPConn
// 	if index != -1 && net_platform.Connected_nodes[index].conn != nil {
// 		conn = net_platform.Connected_nodes[index].conn
// 	} else {
// 		dst_addr, _ := net.ResolveTCPAddr("tcp", addr)
// 		src_addr := net_platform.Self_node.Socket
// 		conn, err = net.DialTCP("tcp", src_addr, dst_addr)

// 		if err != nil {
// 			fmt.Printf("Connection Failed, for node with address: %s\nError: %s\n", addr, err)
// 			net_platform.AddToPreviousNodes(addr)
// 			net_platform.RemoveNodeWithAddress(addr)
// 			if net_platform.node_failure_event_handler != nil {
// 				net_platform.node_failure_event_handler(net_platform, addr)
// 			}
// 			log.Println(err)
// 			return err
// 		}

// 		defer conn.Close()
// 	}
// 	_, err = conn.Write(data) // write into connection i.e send data

// 	if err != nil {
// 		if err == io.EOF {
// 			log.Println("Client ", conn.RemoteAddr(), " disconnected")
// 			conn.Close()
// 			return nil
// 		} else {
// 			log.Println("Failed writing bytes to conn: ", conn, " with error ", err)
// 			conn.Close()
// 			return err
// 		}
// 	}

// 	return nil
// }

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
