package nodes

import (
	"GuthiNetwork/events"
	"fmt"
	"net"
	"os"
)

// So each node should work as server and client simultaneously

type NetworkNode struct {
	NodeID int
	Name   string
	// TCP Addr is akin to socket. So, its only used when its listening for connection, right?
	Socket *net.TCPAddr // TCP address of the current node
}

// Test code
func CheckTCPConnection() {
	tcp_info := "localhost:8080"
	tcp_addr, err := net.ResolveTCPAddr("tcp", tcp_info)
	if err != nil {
		println("Cannot resolve TCP addr info : ", err)
		os.Exit(1)
	}

	socket, err := net.DialTCP("tcp", nil, tcp_addr)

	if err != nil {
		println("Failed to dial tcp server", err.Error())
		os.Exit(1)
	}
	_, err = socket.Write([]byte("Send me time bro"))

	time := make([]byte, 2048)
	// Reading from the server
	for {
		count, err := socket.Read(time)
		if err != nil {
			println("Failed to read from the server")
			break
		}
		println("Time received is : ", string(time[0:count]))
	}
	socket.Close()
}

func IntiateTCPConnection(node *NetworkNode) *net.TCPConn {
	tcp_con, err := net.DialTCP("tcp", nil, node.Socket)
	if err != nil {
		fmt.Println("Failed to initiate tcp connection with the host : ", node)
		return nil
	}
	return tcp_con
}

// Implement a event queue here
type Queue struct {
	capacity uint
	len      uint
	data     []events.Events
}

// To be implemented later on
